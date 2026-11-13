package joiner_runtime

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/pavlo67/common/common/joiner"
)

func TestInterface(t *testing.T) {
	joinerOp := New(nil, nil)

	const keyA1 joiner.InterfaceKey = "KeyA1"
	structA1 := &StructA{}

	const keyA2 joiner.InterfaceKey = "KeyA2"
	structA2 := &StructA{}

	joinerOp.Join(structA1, keyA1)
	joinerOp.Join(structA2, keyA2)

	structA1Joined, ok := joinerOp.Interface(keyA1).(InterfaceA)
	require.True(t, ok)
	require.Equal(t, structA1, structA1Joined)
}

//func TestComponentsAll(t *testing.T) {
//	joiner := New()
//
//	const textA1 = "StructA.TypeKey()"
//	const keyA1 HandlerKey = "KeyA1"
//	structA1 := &StructA{text: textA1}
//	structA3 := &StructA{text: textA1}
//
//	const keyA2 HandlerKey = "KeyA2"
//	structA2 := &StructA{}
//
//	const keyB1 HandlerKey = "KeyB1"
//	structB1 := &StructB{}
//
//	joiner.Join(structA1, keyA1)
//	joiner.Join(structA3, keyA1)
//	joiner.Join(structB1, keyB1)
//	joiner.Join(structA2, keyA2)
//
//	routes := joiner.ComponentsAll(keyA1)
//	require.Equal(t, 2, len(routes))
//
//	for _, component := range routes {
//		require.Equal(t, keyA1, component.Key)
//
//		interfaceA, ok := component.Interface.(InterfaceA)
//		require.True(t, ok)
//		require.NotNil(t, interfaceA)
//
//		text := interfaceA.ActionA()
//		require.Equal(t, textA1, text)
//	}
//
//	require.Equal(t, 1, structA1.NumActionA)
//	require.Equal(t, 1, structA3.NumActionA)
//	require.Equal(t, 0, structA2.NumActionA)
//}

func TestComponentsAllWithSignature(t *testing.T) {
	joinerOp := New(nil, nil)

	const textA1 = "StructA.TypeKey()"
	const keyA1 joiner.InterfaceKey = "KeyA1"
	structA1 := &StructA{text: textA1}

	const keyA2 joiner.InterfaceKey = "KeyA2"
	structA2 := &StructA{text: textA1}

	const keyA3 joiner.InterfaceKey = "KeyA3"
	structA3 := &StructA{text: textA1}

	const keyB1 joiner.InterfaceKey = "KeyB1"
	structB1 := &StructB{}

	joinerOp.Join(structA1, keyA1)
	joinerOp.Join(structA3, keyA3)
	joinerOp.Join(structB1, keyB1)
	joinerOp.Join(structA2, keyA2)

	components := joinerOp.InterfacesAll((*InterfaceA)(nil))
	require.Equal(t, 3, len(components))

	for _, component := range components {
		interfaceA, ok := component.Interface.(InterfaceA)
		require.True(t, ok)
		require.NotNil(t, interfaceA)

		text := interfaceA.ActionA()
		require.Equal(t, textA1, text)
	}

	require.Equal(t, 1, structA1.NumActionA)
	require.Equal(t, 1, structA3.NumActionA)
	require.Equal(t, 1, structA2.NumActionA)
}

func TestCloseAll(t *testing.T) {
	joinerOp := New(nil, nil)

	const textA1 = "StructA.TypeKey()"
	const keyA1 joiner.InterfaceKey = "KeyA1"
	structA1 := &StructA{text: textA1}

	const keyA2 joiner.InterfaceKey = "KeyA2"
	structA2 := &StructA{text: textA1}

	const keyA3 joiner.InterfaceKey = "KeyA3"
	structA3 := &StructA{text: textA1}

	const keyB1 joiner.InterfaceKey = "KeyB1"
	structB1 := &StructB{}

	joinerOp.Join(structA1, keyA1)
	joinerOp.Join(structA3, keyA3)
	joinerOp.Join(structB1, keyB1)
	joinerOp.Join(structA2, keyA2)

	joinerOp.CloseAll()

	require.Equal(t, 1, structA1.NumClose)
	require.Equal(t, 1, structA2.NumClose)
	require.Equal(t, 1, structA3.NumClose)
	require.Equal(t, 1, structB1.NumClose)
}

// InterfaceA (includes Closer) -----------------------------------------------------------------------------------------------------

type InterfaceA interface {
	ActionA() string
}

type StructA struct {
	NumActionA, NumClose int
	text                 string
}

var _ InterfaceA = &StructA{}
var _ joiner.Closer = &StructA{}

func (s *StructA) ActionA() string {
	s.NumActionA++
	fmt.Println("StructA.TypeKey()")
	return s.text
}

func (s *StructA) Close() error {
	s.NumClose++
	fmt.Println("StructA.Close()")
	return nil
}

// InterfaceB (includes Closer) -----------------------------------------------------------------------------------------------------

type InterfaceB interface {
	ActionB() string
}

type StructB struct {
	NumActionB, NumClose int
}

var _ InterfaceB = &StructB{}
var _ joiner.Closer = &StructB{}

func (s *StructB) ActionB() string {
	s.NumActionB++
	fmt.Println("StructB.TypeKey()")
	return "StructB.TypeKey()"
}

func (s *StructB) Close() error {
	s.NumClose++
	fmt.Println("StructB.Close()")
	return nil
}
