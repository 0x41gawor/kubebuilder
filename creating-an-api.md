# Kubernets API
As you might know for the purpose of  managing objects in K8s cluster, `kube-api-server` exposes an REST API.

E.g. for managing a resources of a [kind Pod we have endpoints](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.26/#pod-v1-core):
- `POST /api/v1/namespaces/{namespace}/pods`                -> Creation
- `PATCH /api/v1/namespaces/{namespace}/pods/{name}`        -> Partial Update
- `PUT /api/v1/namespaces/{namespace}/pods/{name}`          -> Replacement (Update) 
- `DELETE /api/v1/namespaces/{namespace}/pods/{name}`       -> Deletion
- `GET /api/v1/namespaces/{namespace}/pods/{name}`          -> Read
- `GET /api/v1/namespaces/{namespace}/pods`                 -> List
- `GET /api/v1/pods`                                        -> List all namespaces
- `GET /api/v1/watch/namespaces/{namespace}/pods/{name}`    -> Watch
- `GET /api/v1/watch/namespaces/{namespace}/pods`           -> Watch list
- `GET /api/v1/watch/pods`                                  -> Watch list all namespaces
- `PATCH /api/v1/namespaces/{namespace}/pods/{name}/status` -> Partially update status
- `GET /api/v1/namespaces/{namespace}/pods/{name}/status`   -> Read status
- `PUT /api/v1/namespaces/{namespace}/pods/{name}/status`   -> Read status

During a https://github.com/0x41gawor/pdmgr/blob/master/crd-playground-lab.md an API just like the one above was created for `Loop` kind.

## Building own API

Kubebuilder is not only a tool for developing operators. It is also a tool for developing own custom resources, so it also has a part of what can be done by defining a CRD and applying it. Kubebuilder additionally adds an opportunity to code the controller and webhooks of such resource. I.e. Kubebuilder helps you scaffold and build CRDs along with controller and webhooks to manage them.

When you run a command like
```sh
kubebuilder create api --group <group> --version <version> --kind <kind>
```

Kubebuilder provides you with GO struct to define the **spec** and **status** of your custom resource. 

**Schema**<br>
This struct serves as the schema for your CRD, specifying which fields are required and their types. Kubebuilder uses annotations to define OpenAPI validations, which get translated to Kubernetes CRD definitions.

**Registration**<br>
The generated CRD is registered with Kubernetes so that it knows about the new resource. Once registered, the custom resource type can be created and managed like any built-in resource (e.g., using kubectl commands).

## GVK 101
Yeah, but what are the group, version and kind in the commmand above?

Let's take a look at endpoints of kind Pod and other built-in kinds: 
```sh
Pods        -> /api/v1/namespaces/{namespace}/pods
Deployments -> /apis/apps/v1/namespaces/{namespace}/deployments
Jobs        -> /apis/batch/v1/namespaces/{namespace}/jobs
Ingresses   -> /apis/networking.k8s.io/v1/namespaces/{namespace}/ingresses
ReplicaSets -> /apis/apps/v1/namespaces/{namespace}/replicasets
```

As you can see the beggining of the enpoints can vary. This is the group.

In Kubernetes, resources are grouped into API groups to organize and version the API endpoints.

Here are some common Kubernetes API groups and examples of resources included in each group:
- `core` (also known as "legacy group")
    - `/api/v1`
    - Group for most fundamental resources
    - Pod, Service, Node, ConfigMap, Secret
- `apps`
    - `apis/apps/v1`
    - higher-level controllers for deploying applications
    - Deployment, ReplicaSet, StatefulSet, DeamonSet
- `batch`
    - `apis/batch/v1`
    - Deals with batch processing and task execution
    - Job, CronJob
- `autoscaling`
    - `apis/autoscaling/v1`
    - Manages automatic scaling
    - HorizontalPodAutoscaler
- `networking`
    - `/apis/networking.k8s.io/v1`
    - Manages networking and connectivity
    - Ingress, NetworkPolicy

Complete list of Api groups can be found here: https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.26/#-strong-api-groups-strong-

As you can see each group has its version at the end. If the API will have some significant changes its version can be raised in the future. But previous version still can be used.

For you, as a developer of Custom Resource designing an groups and version are very helpful.
### Group
When you are designing an api and want to group some of its part because of related functionality, you create a group. In here, a group is a set of related Kinds.
### Version
When you are designing an api and want to change it over time you need to introduce versioning, so the client is aware of it and experience only expected behavior. In here, what you are versioning is a group.

#### Versioning conventions
In Kubernetes each group is versioned to allow for API stability and evolution. The version can be either:

- Alpha (v1alpha1, v1alpha2, etc.): Indicates that the feature is new, may not be complete, and is subject to breaking changes without notice.
- Beta (v1beta1, v1beta2, etc.): More stable than alpha, indicating that the feature is well-tested but still subject to potential changes.
- Stable (v1, v2, etc.): Indicates a mature and stable feature that is unlikely to change.

The structure of the API groups and their versioning allows Kubernetes to evolve over time, adding new features and refining existing ones while maintaining compatibility with existing deployments.

### Object vs Resource vs Kind

#### Object

[Object or ApiObject](https://kubernetes.io/docs/concepts/overview/working-with-objects/) is a key concepts of how Kubernets works. Cluster state is described throught them.

Kubernetes objects are persistent entities in the Kubernetes system. Kubernetes uses these entities to represent the state of your cluster. Specifically, they can describe:
- What containerized applications are running (and on which nodes)
- The resources available to those applications
- The policies around how those applications behave, such as restart policies, upgrades, and fault-tolerance

A Kubernetes object is a "record of intent"--once you create the object, the Kubernetes system will constantly work to ensure that the object exists.

Almost every Kubernetes object includes two nested object fields that govern the object's configuration: the object spec and the object status. For objects that have a spec, you have to set this when you create the object, providing a description of the characteristics you want the resource to have: its desired state.

The status describes the current state of the object, supplied and updated by the Kubernetes system and its components. The Kubernetes control plane continually and actively manages every object's actual state to match the desired state you supplied.

#### Kind

ApiObject can be of some Kind. It is a way to organize objects of the same set of attributes.

Different kind of objects have different set of status fields.

A kind is a type that defines the schema and properties of a Kubernetes object.

#### Resource

This term is more related to YOU. Actual user of objects. The objects that you will create are your resources.

A resource is an API endpoint in the Kubernetes API that refers to a group of objects of the same kind. In other words, it represents all instances of a specific kind in the cluster.

#### Summary

Exemplary kind is: Pod, Deployment, Service <br>
Exemplary object is an actual instance created by kubernetes after hitting some POST endpoints such as `POST /api/v1/namespaces/{namespace}/pods`. Objects are stored in etcd.
Exemplary resource are pods, deployments, services, or custom resoruces like widget aftrer defining a CRD. 

So kind in an abstract term defined by Kubernetes. They defined that kind "Pod" has this schema etc... The resource are your instances of this kind in you k8s cluster. You can manage this resource with endpoint `api/v1/pods`. An object is instance of given kind created after `POST /api/v1/namespaces/{namespace}/pods`.
