package platform

import (
	"fmt"

	"github.com/AlexsJones/gravitywell/state"
	"github.com/fatih/color"
	v1polbeta "k8s.io/api/policy/v1beta1"
	v12 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	_ "k8s.io/client-go/plugin/pkg/client/auth/gcp"
)

func execPodDisruptionBudgetResouce(k kubernetes.Interface, pdb *v1polbeta.PodDisruptionBudget, namespace string, dryRun bool, tryUpdate bool) (state.State, error) {
	color.Blue("Found PodDisruptionBudget resource")
	pdbclient := k.PolicyV1beta1().PodDisruptionBudgets(namespace)

	if dryRun {
		_, err := pdbclient.Get(pdb.Name, v12.GetOptions{})
		if err != nil {
			color.Red(fmt.Sprintf("DRY-RUN: PodDisruptionBudget resource %s does not exist\n", pdb.Name))
			return state.EDeploymentStateNotExists, err
		} else {
			color.Blue(fmt.Sprintf("DRY-RUN: PodDisruptionBudget resource %s exists\n", pdb.Name))
			return state.EDeploymentStateExists, nil
		}
	}

	_, err := pdbclient.Create(pdb)
	if err != nil {
		if !tryUpdate {
			color.Cyan("PodDisruptionBudget already exists - Cowardly refusing to overwrite")
			return state.EDeploymentStateExists, err
		}
		_, err := pdbclient.Update(pdb)
		if err != nil {
			color.Red("PodDisruptionBudget could not be updated")
			return state.EDeploymentStateCantUpdate, err
		}
		color.Blue("Configmap updated")
	}
	return state.EDeploymentStateOkay, nil
}