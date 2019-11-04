package main

import (
	"context"
	"github.com/docker/docker/api/types"
	client2 "github.com/docker/docker/client"
	"time"
)

var client *client2.Client

func InitClient() {
	var err error
	client, err = client2.NewEnvClient()
	if err != nil {
		panic(err)
	}
}

func tryInitClient() (r bool) {
	defer func() {
		if e := recover(); e != nil {
			r = false
		}
	}()
	InitClient()
	r = true
	return
}

func containers(all bool) (c []types.Container) {
	c, _ = client.ContainerList(context.Background(), types.ContainerListOptions{
		All: all,
	})
	return
}

func startContainer(id string) bool {
	for _, c := range containers(true) {
		if c.ID == id {
			return client.ContainerStart(context.Background(), c.ID, types.ContainerStartOptions{}) == nil
		}
	}
	return false
}

func startContainerByName(name string) bool {
	for _, c := range containers(true) {
		if sliceContains(c.Names, name) {
			return client.ContainerStart(context.Background(), c.ID, types.ContainerStartOptions{}) == nil
		}
	}
	return false
}

func stopContainer(id string) bool {
	tm := time.Millisecond * 10
	for _, c := range containers(true) {
		if c.ID == id {
			return client.ContainerStop(context.Background(), c.ID, &tm) == nil
		}
	}
	return false
}

func stopContainerByName(name string) bool {
	tm := time.Millisecond * 10
	for _, c := range containers(true) {
		if sliceContains(c.Names, name) {
			return client.ContainerStop(context.Background(), c.ID, &tm) == nil
		}
	}
	return false
}
func sliceContains(str []string, s string) bool {
	for _, x := range str {
		if x == s {
			return true
		}
	}
	return false
}
