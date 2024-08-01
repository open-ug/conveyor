package v1

import "k8s.io/apimachinery/pkg/runtime"

// DeepCopyInto copies all properties of this object into another object of the
// same type that is provided as a pointer.
func (in *Application) DeepCopyInto(out *Application) {
	out.TypeMeta = in.TypeMeta
	out.ObjectMeta = in.ObjectMeta
	out.Spec = ApplicationSpec{
		AppName: in.Spec.AppName,
		Volumes: in.Spec.Volumes,
		Ports:   in.Spec.Ports,
		EnvFrom: in.Spec.EnvFrom,
		Image:   in.Spec.Image,
	}
}

// DeepCopyObject returns a generically typed copy of an object
func (in *Application) DeepCopyObject() runtime.Object {
	out := Application{}
	in.DeepCopyInto(&out)

	return &out
}

// DeepCopyObject returns a generically typed copy of an object
func (in *ApplicationList) DeepCopyObject() runtime.Object {
	out := ApplicationList{}
	out.TypeMeta = in.TypeMeta
	out.ListMeta = in.ListMeta

	if in.Items != nil {
		out.Items = make([]Application, len(in.Items))
		for i := range in.Items {
			in.Items[i].DeepCopyInto(&out.Items[i])
		}
	}

	return &out
}
