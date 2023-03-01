harvester-event-logger
========

## Background
In harvester, the `rancher-logging` is deployed to have the full logging feature. 

An resource of [eventtailer](https://github.com/rancher/charts/blob/dev-v2.7/charts/rancher-logging/100.1.1%2Bup3.17.3/templates/clusterrole.yaml#L201) will lead to the deployment of an event router.

Harvester defines an [eventtailer](https://github.com/harvester/harvester-installer/blob/cc6e9265942ac08506de20e3229d6844021c41cb/pkg/config/templates/rancherd-15-logging.yaml#L61) deployed, which is using a default image `banzaicloud/eventrouter:v0.1.0`. The project and its upsteam are not in active developing status. [banzaicloud/eventrouter](https://github.com/banzaicloud/eventrouter/commits/master)

To solve https://github.com/harvester/security/issues/19, we create this project as an replacement.

note:

The clusterrole is defined in:
[rancher-logging](https://github.com/rancher/charts/blob/dev-v2.7/charts/rancher-logging/100.1.1%2Bup3.17.3/templates/clusterrole.yaml#L82)

A separate clusterrole is not defined in this project.

## Building

When you do not find `vendor` in the source code, then run

`go mod vendor` to download the vendor.


`make`

will build, test and package the image

check the image via

`docker image ls "rancher/harvester-event-logger:dev"`

the output will be like:

```
REPOSITORY                       TAG       IMAGE ID       CREATED       SIZE
rancher/harvester-event-logger   dev       9b7f77a996f0   3 hours ago   78.9MB
```

## Deployment

### Manual local test

From building PC:
```
docker save -o hel.img rancher/harvester-event-logger:dev
scp hel.img rancher@192.168.122.206://home/rancher
```

From Harvester cluster, ssh into NODE:
```
sudo -i

docker image load -i /home/rancher/hel.img

kubectl set image pod -n cattle-logging-system harvester-default-event-tailer-0 *=rancher/harvester-event-logger:dev
```

### Auto test

As most Harvester projects, upload the image into an repository, pack the image into harvester ISO, change the default image with this one. 

Those will be done after this project is adopted.

## License

Copyright (c) 2023 [Rancher Labs, Inc.](http://rancher.com)

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

[http://www.apache.org/licenses/LICENSE-2.0](http://www.apache.org/licenses/LICENSE-2.0)

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
