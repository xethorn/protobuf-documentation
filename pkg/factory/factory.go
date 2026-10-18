package factory

import (
	"fmt"
	"math/rand"
	"time"

	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/descriptorpb"
)

type DescriptorFactory struct {
	rnd *rand.Rand
}

func NewDescriptorFactory(seed int64) *DescriptorFactory {
	if seed == 0 {
		seed = time.Now().UnixNano()
	}
	return &DescriptorFactory{
		rnd: rand.New(rand.NewSource(seed)),
	}
}

func (f *DescriptorFactory) RandomString(prefix string) string {
	return fmt.Sprintf("%s_%d", prefix, f.rnd.Intn(100000))
}

func (f *DescriptorFactory) CreateFileDescriptor(pkg string, serviceCount, messageCount, enumCount int) *descriptorpb.FileDescriptorProto {
	fileName := f.RandomString("file") + ".proto"
	fileProto := &descriptorpb.FileDescriptorProto{
		Name:    proto.String(fileName),
		Package: proto.String(pkg),
		Syntax:  proto.String("proto3"),
		SourceCodeInfo: &descriptorpb.SourceCodeInfo{
			Location: []*descriptorpb.SourceCodeInfo_Location{
				{
					Path:            []int32{},
					LeadingComments: proto.String("File level documentation comment for " + fileName),
				},
			},
		},
	}

	for i := 0; i < messageCount; i++ {
		msgName := fmt.Sprintf("Message%d", i)
		msgPath := []int32{4, int32(i)}
		fileProto.SourceCodeInfo.Location = append(fileProto.SourceCodeInfo.Location, &descriptorpb.SourceCodeInfo_Location{
			Path:            msgPath,
			LeadingComments: proto.String("Documentation for " + msgName),
		})

		msg := &descriptorpb.DescriptorProto{
			Name: proto.String(msgName),
		}

		// Add fields to message
		for j := 0; j < 3; j++ {
			fldName := fmt.Sprintf("field_%d", j)
			fldPath := append([]int32{}, msgPath...)
			fldPath = append(fldPath, 2, int32(j))
			fileProto.SourceCodeInfo.Location = append(fileProto.SourceCodeInfo.Location, &descriptorpb.SourceCodeInfo_Location{
				Path:            fldPath,
				LeadingComments: proto.String("Field documentation for " + fldName),
			})

			msg.Field = append(msg.Field, &descriptorpb.FieldDescriptorProto{
				Name:     proto.String(fldName),
				Number:   proto.Int32(int32(j + 1)),
				Type:     descriptorpb.FieldDescriptorProto_TYPE_STRING.Enum(),
				Label:    descriptorpb.FieldDescriptorProto_LABEL_OPTIONAL.Enum(),
				JsonName: proto.String(fldName),
			})
		}

		fileProto.MessageType = append(fileProto.MessageType, msg)
	}

	for i := 0; i < enumCount; i++ {
		enumName := fmt.Sprintf("Enum%d", i)
		enumPath := []int32{5, int32(i)}
		fileProto.SourceCodeInfo.Location = append(fileProto.SourceCodeInfo.Location, &descriptorpb.SourceCodeInfo_Location{
			Path:            enumPath,
			LeadingComments: proto.String("Documentation for " + enumName),
		})

		enumProto := &descriptorpb.EnumDescriptorProto{
			Name: proto.String(enumName),
			Value: []*descriptorpb.EnumValueDescriptorProto{
				{
					Name:   proto.String(enumName + "_UNSPECIFIED"),
					Number: proto.Int32(0),
				},
				{
					Name:   proto.String(enumName + "_ACTIVE"),
					Number: proto.Int32(1),
				},
			},
		}

		fileProto.EnumType = append(fileProto.EnumType, enumProto)
	}

	for i := 0; i < serviceCount; i++ {
		svcName := fmt.Sprintf("Service%d", i)
		svcPath := []int32{6, int32(i)}
		fileProto.SourceCodeInfo.Location = append(fileProto.SourceCodeInfo.Location, &descriptorpb.SourceCodeInfo_Location{
			Path:            svcPath,
			LeadingComments: proto.String("Service documentation for " + svcName),
		})

		svcProto := &descriptorpb.ServiceDescriptorProto{
			Name: proto.String(svcName),
		}

		// Add method
		methodPath := append([]int32{}, svcPath...)
		methodPath = append(methodPath, 2, 0)
		fileProto.SourceCodeInfo.Location = append(fileProto.SourceCodeInfo.Location, &descriptorpb.SourceCodeInfo_Location{
			Path:            methodPath,
			LeadingComments: proto.String("RPC method comment"),
		})

		svcProto.Method = append(svcProto.Method, &descriptorpb.MethodDescriptorProto{
			Name:       proto.String("ExecuteMethod"),
			InputType:  proto.String("." + pkg + ".Message0"),
			OutputType: proto.String("." + pkg + ".Message0"),
		})

		fileProto.Service = append(fileProto.Service, svcProto)
	}

	return fileProto
}
