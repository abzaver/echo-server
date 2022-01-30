#!/bin/sh

kubectl exec deployments/svc-hello -n svc-hello-ns -- /bin/busybox killall -HUP svc-hello