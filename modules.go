package main

import (
	_ "github.com/contentanalyst/registrator/consul"
	_ "github.com/contentanalyst/registrator/consulkv"
	_ "github.com/contentanalyst/registrator/etcd"
	_ "github.com/contentanalyst/registrator/rancher"
	_ "github.com/contentanalyst/registrator/skydns2"
	_ "github.com/contentanalyst/registrator/zookeeper"
)
