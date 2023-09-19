/*

	MIT License

	Copyright (c) Microsoft Corporation.

	Permission is hereby granted, free of charge, to any person obtaining a copy
	of this software and associated documentation files (the "Software"), to deal
	in the Software without restriction, including without limitation the rights
	to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
	copies of the Software, and to permit persons to whom the Software is
	furnished to do so, subject to the following conditions:

	The above copyright notice and this permission notice shall be included in all
	copies or substantial portions of the Software.

	THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
	IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
	FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
	AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
	LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
	OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
	SOFTWARE

*/

package v1

import (
	"context"
	"strings"

	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/util/validation/field"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/webhook"

	configv1 "gopls-workspace/apis/config/v1"
	configutils "gopls-workspace/configutils"

	apierrors "k8s.io/apimachinery/pkg/api/errors"
)

// log is for logging in this package.
var modellog = logf.Log.WithName("model-resource")
var myModelClient client.Client
var modelValidationPolicies []configv1.ValidationPolicy

func (r *Model) SetupWebhookWithManager(mgr ctrl.Manager) error {
	myModelClient = mgr.GetClient()
	mgr.GetFieldIndexer().IndexField(context.Background(), &Model{}, ".spec.displayName", func(rawObj client.Object) []string {
		model := rawObj.(*Model)
		return []string{model.Spec.DisplayName}
	})

	dict, _ := configutils.GetValidationPoilicies()
	if v, ok := dict["model"]; ok {
		modelValidationPolicies = v
	}

	return ctrl.NewWebhookManagedBy(mgr).
		For(r).
		Complete()
}

// TODO(user): EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!

var _ webhook.Defaulter = &Model{}

// Default implements webhook.Defaulter so a webhook will be registered for the type
func (r *Model) Default() {
	modellog.Info("default", "name", r.Name)

	if r.Spec.DisplayName == "" {
		r.Spec.DisplayName = r.ObjectMeta.Name
	}
}

// TODO(user): change verbs to "verbs=create;update;delete" if you want to enable deletion validation.

var _ webhook.Validator = &Model{}

// ValidateCreate implements webhook.Validator so a webhook will be registered for the type
func (r *Model) ValidateCreate() error {
	modellog.Info("validate create", "name", r.Name)

	return r.validateCreateModel()
}

// ValidateUpdate implements webhook.Validator so a webhook will be registered for the type
func (r *Model) ValidateUpdate(old runtime.Object) error {
	modellog.Info("validate update", "name", r.Name)

	return r.validateUpdateModel()
}

// ValidateDelete implements webhook.Validator so a webhook will be registered for the type
func (r *Model) ValidateDelete() error {
	modellog.Info("validate delete", "name", r.Name)

	// TODO(user): fill in your validation logic upon object deletion.
	return nil
}

func (r *Model) validateCreateModel() error {
	var allErrs field.ErrorList
	var models ModelList
	err := myModelClient.List(context.Background(), &models, client.InNamespace(r.Namespace), client.MatchingFields{".spec.displayName": r.Spec.DisplayName})
	if err != nil {
		allErrs = append(allErrs, field.InternalError(&field.Path{}, err))
		return apierrors.NewInvalid(schema.GroupKind{Group: "ai.symphony", Kind: "Model"}, r.Name, allErrs)
	}
	if len(models.Items) != 0 {
		allErrs = append(allErrs, field.Invalid(field.NewPath("spec").Child("displayName"), r.Spec.DisplayName, "model display name is already taken"))
		return apierrors.NewInvalid(schema.GroupKind{Group: "ai.symphony", Kind: "Model"}, r.Name, allErrs)
	}
	if len(modelValidationPolicies) > 0 {
		err := myModelClient.List(context.Background(), &models, client.InNamespace(r.Namespace), &client.ListOptions{})
		if err != nil {
			allErrs = append(allErrs, field.InternalError(&field.Path{}, err))
			return apierrors.NewInvalid(schema.GroupKind{Group: "ai.symphony", Kind: "Model"}, r.Name, allErrs)
		}
		for _, p := range modelValidationPolicies {
			pack := extractModelValidationPack(models, p)
			ret, err := configutils.CheckValidationPack(r.ObjectMeta.Name, readModelValiationTarget(r, p), p.ValidationType, pack)
			if err != nil {
				return err
			}
			if ret != "" {
				allErrs = append(allErrs, field.Forbidden(&field.Path{}, strings.ReplaceAll(p.Message, "%s", ret)))
				return apierrors.NewInvalid(schema.GroupKind{Group: "ai.symphony", Kind: "Model"}, r.Name, allErrs)
			}
		}
	}
	return nil
}

func (r *Model) validateUpdateModel() error {
	var allErrs field.ErrorList
	var models ModelList
	err := myModelClient.List(context.Background(), &models, client.InNamespace(r.Namespace), client.MatchingFields{".spec.displayName": r.Spec.DisplayName})
	if err != nil {
		allErrs = append(allErrs, field.InternalError(&field.Path{}, err))
		return apierrors.NewInvalid(schema.GroupKind{Group: "ai.symphony", Kind: "Model"}, r.Name, allErrs)
	}
	if !(len(models.Items) == 0 || len(models.Items) == 1 && models.Items[0].ObjectMeta.Name == r.ObjectMeta.Name) {
		allErrs = append(allErrs, field.Invalid(field.NewPath("spec").Child("displayName"), r.Spec.DisplayName, "model display name is already taken"))
		return apierrors.NewInvalid(schema.GroupKind{Group: "ai.symphony", Kind: "Model"}, r.Name, allErrs)
	}
	if len(modelValidationPolicies) > 0 {
		err := myModelClient.List(context.Background(), &models, client.InNamespace(r.Namespace), &client.ListOptions{})
		if err != nil {
			allErrs = append(allErrs, field.InternalError(&field.Path{}, err))
			return apierrors.NewInvalid(schema.GroupKind{Group: "ai.symphony", Kind: "Model"}, r.Name, allErrs)
		}
		for _, p := range modelValidationPolicies {
			pack := extractModelValidationPack(models, p)
			ret, err := configutils.CheckValidationPack(r.ObjectMeta.Name, readModelValiationTarget(r, p), p.ValidationType, pack)
			if err != nil {
				return err
			}
			if ret != "" {
				allErrs = append(allErrs, field.Forbidden(&field.Path{}, strings.ReplaceAll(p.Message, "%s", ret)))
				return apierrors.NewInvalid(schema.GroupKind{Group: "ai.symphony", Kind: "Model"}, r.Name, allErrs)
			}
		}
	}
	return nil
}

func readModelValiationTarget(model *Model, p configv1.ValidationPolicy) string {
	if p.SelectorType == "properties" {
		if v, ok := model.Spec.Properties[p.SpecField]; ok {
			return v
		}
	}
	return ""
}

func extractModelValidationPack(list ModelList, p configv1.ValidationPolicy) []configv1.ValidationStruct {
	pack := make([]configv1.ValidationStruct, 0)
	for _, t := range list.Items {
		s := configv1.ValidationStruct{}
		if p.SelectorType == "properties" {
			if v, ok := t.Spec.Properties[p.SpecField]; ok {
				s.Field = v
				s.Name = t.ObjectMeta.Name
				pack = append(pack, s)
			}
		}
	}
	return pack
}
