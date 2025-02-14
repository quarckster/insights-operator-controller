/*
Copyright © 2019 Red Hat, Inc.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package server

import (
	"encoding/json"
	"github.com/redhatinsighs/insights-operator-controller/storage"
	"io/ioutil"
	"net/http"
	"testing"
)

const API_URL = "http://localhost:8080/api/v1/"

func readListOfClusters(t *testing.T) []storage.Cluster {
	response, err := http.Get(API_URL + "client/cluster")
	if err != nil {
		t.Errorf("Communication error with the server %v", err)
	}
	if response.StatusCode != http.StatusOK {
		t.Errorf("Expected HTTP status 200 OK, got %d", response.StatusCode)
	}

	body, readErr := ioutil.ReadAll(response.Body)
	defer response.Body.Close()

	if readErr != nil {
		t.Errorf("Unable to read response body")
	}

	clusters := []storage.Cluster{}
	err = json.Unmarshal(body, &clusters)
	if err != nil {
		t.Error(err)
	}
	return clusters
}

func createClusterTestStep(t *testing.T, clusterId string, clusterName string) {
	var client http.Client

	url := API_URL + "client/cluster/" + clusterId + "/" + clusterName
	request, err := http.NewRequest("POST", url, nil)
	if err != nil {
		t.Errorf("Error creating request %v", err)
	}

	response, err := client.Do(request)
	if err != nil {
		t.Errorf("Communication error with the server %v", err)
	}
	if response.StatusCode != http.StatusCreated {
		t.Errorf("Expected HTTP status 201 Created, got %d", response.StatusCode)
	}
}

// Test step that try to delete the cluster via REST API.
func deleteClusterTestStep(t *testing.T, clusterId string) {
	var client http.Client

	url := API_URL + "client/cluster/" + clusterId
	request, err := http.NewRequest("DELETE", url, nil)
	if err != nil {
		t.Errorf("Error creating request %v", err)
	}

	response, err := client.Do(request)
	if err != nil {
		t.Errorf("Communication error with the server %v", err)
	}
	if response.StatusCode != http.StatusAccepted {
		t.Errorf("Expected HTTP status 202 Accepted, got %d", response.StatusCode)
	}
}

func compareClusters(t *testing.T, clusters []storage.Cluster, expected []storage.Cluster) {
	if len(clusters) != len(expected) {
		t.Errorf("%d clusters are expected, but got %d", len(expected), len(clusters))
	}

	for i := 0; i < len(expected); i++ {
		if clusters[i] != expected[i] {
			t.Errorf("Different cluster info returned: %v != %v", clusters[i], expected[i])
		}
	}
}

func TestGetListOfClusters(t *testing.T) {
	clusters := readListOfClusters(t)

	expected := []storage.Cluster{
		{0, "cluster0"},
		{1, "cluster1"},
		{2, "cluster2"},
		{3, "cluster3"},
		{4, "cluster4"},
	}
	compareClusters(t, clusters, expected)
}

func TestAddCluster(t *testing.T) {
	var client http.Client

	request, err := http.NewRequest("POST", API_URL+"client/cluster/5/cluster5", nil)
	if err != nil {
		t.Errorf("Error creating request %v", err)
	}

	response, err := client.Do(request)
	if response.StatusCode != http.StatusCreated {
		t.Errorf("Expected HTTP status 201 Created, got %d", response.StatusCode)
	}

	clusters := readListOfClusters(t)

	expected := []storage.Cluster{
		{0, "cluster0"},
		{1, "cluster1"},
		{2, "cluster2"},
		{3, "cluster3"},
		{4, "cluster4"},
		{5, "cluster5"},
	}
	compareClusters(t, clusters, expected)
}

// Check that cluster can be deleted via REST API.
func TestDeleteCluster(t *testing.T) {
	deleteClusterTestStep(t, "5")

	clusters := readListOfClusters(t)

	expected := []storage.Cluster{
		{0, "cluster0"},
		{1, "cluster1"},
		{2, "cluster2"},
		{3, "cluster3"},
		{4, "cluster4"},
	}
	compareClusters(t, clusters, expected)
}

// Check that another cluster can be deleted via REST API.
func TestDeleteAnotherCluster(t *testing.T) {
	deleteClusterTestStep(t, "4")

	clusters := readListOfClusters(t)

	expected := []storage.Cluster{
		{0, "cluster0"},
		{1, "cluster1"},
		{2, "cluster2"},
		{3, "cluster3"},
	}
	compareClusters(t, clusters, expected)
}

// Check how is nonexisting cluster handled.
func TestDeleteNonexistentCluster(t *testing.T) {
	deleteClusterTestStep(t, "40")

	clusters := readListOfClusters(t)

	expected := []storage.Cluster{
		{0, "cluster0"},
		{1, "cluster1"},
		{2, "cluster2"},
		{3, "cluster3"},
	}
	compareClusters(t, clusters, expected)
}

// Check the database after all clusters are deleted.
func _TestDeleteAllClusters(t *testing.T) {
	deleteClusterTestStep(t, "0")
	deleteClusterTestStep(t, "1")
	deleteClusterTestStep(t, "2")
	deleteClusterTestStep(t, "3")

	clusters := readListOfClusters(t)

	expected := []storage.Cluster{}
	compareClusters(t, clusters, expected)
}

// Check if new cluster can be created
func TestCreateCluster(t *testing.T) {
	createClusterTestStep(t, "5", "cluster5")

	clusters := readListOfClusters(t)

	expected := []storage.Cluster{
		{0, "cluster0"},
		{1, "cluster1"},
		{2, "cluster2"},
		{3, "cluster3"},
		{5, "cluster5"},
	}
	compareClusters(t, clusters, expected)
}

// Check if new cluster can be created
func TestCreateCluster9(t *testing.T) {
	createClusterTestStep(t, "9", "cluster9")

	clusters := readListOfClusters(t)

	expected := []storage.Cluster{
		{0, "cluster0"},
		{1, "cluster1"},
		{2, "cluster2"},
		{3, "cluster3"},
		{5, "cluster5"},
		{9, "cluster9"},
	}
	compareClusters(t, clusters, expected)
}
