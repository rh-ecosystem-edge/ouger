package codec

import (
	"bytes"
	"fmt"

	"github.com/openshift/api"
	yaml "gopkg.in/yaml.v2"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	serializer "k8s.io/apimachinery/pkg/runtime/serializer"
	jsonserializer "k8s.io/apimachinery/pkg/runtime/serializer/json"
	kapiregistrationv1 "k8s.io/kube-aggregator/pkg/apis/apiregistration/v1"
	"k8s.io/kubectl/pkg/scheme"
)

func Decode(input []byte) ([]byte, error) {
	var buf bytes.Buffer
	if _, ok := tryFindProto(input); ok {
		decoder := scheme.Codecs.UniversalDeserializer()
		encoder := jsonserializer.NewYAMLSerializer(jsonserializer.DefaultMetaFactory, scheme.Scheme, scheme.Scheme)

		obj, _, err := decoder.Decode(input, nil, nil)
		if err != nil {
			return nil, fmt.Errorf("unable to decode input: %s", err)
		}

		err = encoder.Encode(obj, &buf)
		if err != nil {
			return nil, fmt.Errorf("unable to encode input: %s", err)
		}
	} else {
		buf.Write(input)
	}

	return buf.Bytes(), nil
}

func Encode(input []byte) ([]byte, error) {
	typeMeta, err := typeMetaFromYaml(input)
	if err != nil {
		return nil, fmt.Errorf("unable to parse input: %s", err)
	}

	encoder, err := newEncoder(typeMeta)
	if err != nil {
		return nil, fmt.Errorf("unable to get encoder: %s", err)
	}

	decoder := scheme.Codecs.UniversalDeserializer()
	obj, _, err := decoder.Decode(input, nil, nil)
	if err != nil {
		return nil, fmt.Errorf("unable to decode input: %s", err)
	}

	var buf bytes.Buffer
	err = encoder.Encode(obj, &buf)
	if err != nil {
		// just return the raw value
		buf.Write(input)
	}

	return buf.Bytes(), nil
}

func init() {
	api.Install(scheme.Scheme)
	api.InstallKube(scheme.Scheme)

	builder := runtime.NewSchemeBuilder(
		kapiregistrationv1.AddToScheme,
	)
	builder.AddToScheme(scheme.Scheme)
}

var protobufMediaType = "application/vnd.kubernetes.protobuf"

var protoEncodingPrefix = []byte{0x6b, 0x38, 0x73, 0x00}

func tryFindProto(in []byte) ([]byte, bool) {
	i := bytes.Index(in, protoEncodingPrefix)
	if i >= 0 && i < len(in) {
		return in[i:], true
	}
	return nil, false
}


func typeMetaFromYaml(in []byte) (*runtime.TypeMeta, error) {
	var meta runtime.TypeMeta
	yaml.Unmarshal(in, &meta)
	return &meta, nil
}

func newEncoder(typeMeta *runtime.TypeMeta) (runtime.Encoder, error) {
	codecs := serializer.NewCodecFactory(runtime.NewScheme())
	mediaTypes := codecs.SupportedMediaTypes()
	info, ok := runtime.SerializerInfoForMediaType(mediaTypes, protobufMediaType)
	if !ok {
		if len(mediaTypes) == 0 {
			return nil, fmt.Errorf("no serializers registered for %v", mediaTypes)
		}
		info = mediaTypes[0]
	}
	gv, err := schema.ParseGroupVersion(typeMeta.APIVersion)
	if err != nil {
		return nil, fmt.Errorf("unable to parse meta APIVersion '%s': %s", typeMeta.APIVersion, err)
	}
	return scheme.Codecs.EncoderForVersion(info.Serializer, gv), nil
}
