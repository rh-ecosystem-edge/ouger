package main

import (
	"fmt"
	"github.com/openshift/api"
	yaml "gopkg.in/yaml.v2"
	"io/ioutil"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	serializer "k8s.io/apimachinery/pkg/runtime/serializer"
	jsonserializer "k8s.io/apimachinery/pkg/runtime/serializer/json"
	"k8s.io/kubectl/pkg/scheme"
	"os"
)

var ProtobufMediaType = "application/vnd.kubernetes.protobuf"

func main() {
	stdin, err := ioutil.ReadAll(os.Stdin)

	if err != nil {
		panic(fmt.Errorf("unable to read data from stdin: %v", err))
	}

	api.Install(scheme.Scheme)
	api.InstallKube(scheme.Scheme)

	if os.Args[1] == "decode" {
		decoder := scheme.Codecs.UniversalDeserializer()
		encoder := jsonserializer.NewYAMLSerializer(jsonserializer.DefaultMetaFactory, scheme.Scheme, scheme.Scheme)
		obj, _, err := decoder.Decode(stdin, nil, nil)
		if err != nil {
			panic(err)
		}
		err = encoder.Encode(obj, os.Stdout)
		if err != nil {
			panic(err)
		}
	} else if os.Args[1] == "encode" {
		typeMeta, err := typeMetaFromYaml(stdin)
		if err != nil {
			panic(err)
		}

		encoder, err := newEncoder(typeMeta)
		if err != nil {
			panic(err)
		}

		decoder := scheme.Codecs.UniversalDeserializer()
		obj, _, err := decoder.Decode(stdin, nil, nil)
		if err != nil {
			panic(err)
		}

		err = encoder.Encode(obj, os.Stdout)
		if err != nil {
			panic(err)
		}
	} else {
		panic(fmt.Errorf("invalid argument: %v", os.Args[1]))
	}
}

func typeMetaFromYaml(in []byte) (*runtime.TypeMeta, error) {
	var meta runtime.TypeMeta
	yaml.Unmarshal(in, &meta)
	return &meta, nil
}

func newEncoder(typeMeta *runtime.TypeMeta) (runtime.Encoder, error) {
	codecs := serializer.NewCodecFactory(runtime.NewScheme())
	mediaTypes := codecs.SupportedMediaTypes()
	info, ok := runtime.SerializerInfoForMediaType(mediaTypes, ProtobufMediaType)
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
