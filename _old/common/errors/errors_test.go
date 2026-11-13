package errors

import (
	"fmt"
	"log"
	"testing"

	"github.com/pavlo67/common/common"
	"github.com/stretchr/testify/require"

	"github.com/pkg/errors"
)

func TestCommonErrorKey(t *testing.T) {
	testKey1 := Key("test_key1")
	ke1 := CommonError(testKey1)
	require.Equalf(t, testKey1, ke1.Key(), "%#v", ke1)

	testKey2 := Key("test_key2")
	ke2 := CommonError(testKey2, common.Map{"error": "q"})
	require.Equalf(t, testKey2, ke2.Key(), "%#v", ke2)

	testKey3 := Key("test_key3")
	testMap3 := common.Map{"error": "q"}
	ke3 := CommonError(testKey3, testMap3)
	require.Equalf(t, testKey3, ke3.Key(), "%#v", ke3)
	require.Equalf(t, testMap3, ke3.Data(), "%#v", ke3)
}

func TestCommonErrorString(t *testing.T) {
	const TestKey Key = "test"
	err1 := CommonError(TestKey, common.Map{"a": "b"})
	log.Print("ERR1")
	log.Printf("%ss:  %s", "%", err1)
	log.Printf("%sv:  %v", "%", err1)
	log.Printf("%s#v: %#v\n\n", "%", err1)

	err2 := CommonError(err1, "111", errors.New("222"), err1) //
	log.Print("ERR2")
	log.Printf("%ss:  %s", "%", err2)
	log.Printf("%sv:  %v", "%", err2)
	log.Printf("%s#v: %#v", "%", err2)

	errors.New("222")
}

func TestNil(t *testing.T) {
	var err error

	err = CommonError()
	require.Truef(t, err == nil, "%#v", err)

	err = CommonError(nil)
	require.Truef(t, err == nil, "%#v", err)

	err = CommonError(nil, nil)
	require.Truef(t, err == nil, "%#v", err)

	err = CommonError(nil, nil, nil)
	require.Truef(t, err == nil, "%#v", err)

}

func TestWrapf(t *testing.T) {
	err := errors.Wrapf(errors.New("eeeeee"), "22222 %s", "111")
	log.Print(err)
	err1 := CommonError(err, "can't init records.Operator")
	log.Print(fmt.Errorf("error calling .Run() for component (%s): %#v", "name", err1))
	log.Print(fmt.Errorf("error calling .Run() for component (%s): %s", "name", err1))
}

//func TestCommonErrorKey(t *testing.T) {
//	testKey1 := Key("test_key1")
//	ke1 := CommonError(testKey1, nil)
//	require.Equalf(t, testKey1, ke1.Key(), "%#v", ke1)
//
//	testKey2 := Key("test_key2")
//	ke2 := CommonError(testKey2, common.Map{"error": "q"})
//	require.Equalf(t, testKey2, ke2.Key(), "%#v", ke2)
//
//	testKey3 := Key("test_key3")
//	ke3 := CommonError(testKey3, common.Map{"error": "q"})
//	require.Equalf(t, testKey3, ke3.Key(), "%#v", ke3)
//
//	testKey4 := Key("test_key4")
//	ke4 := CommonError(testKey4, common.Map{"error": "q"})
//	require.Equalf(t, testKey4, ke4.Key(), "%#v", ke4)
//
//}
//
//func TestCommonErrorString(t *testing.T) {
//	const TestKey Key = "test"
//	err1 := CommonError(TestKey, common.Map{"a": "b"})
//	log.Print("ERR1")
//	log.Printf("%ss:  %s", "%", err1)
//	log.Printf("%sv:  %v", "%", err1)
//	log.Printf("%s#v: %#v\n\n", "%", err1)
//
//	err2 := CommonError(err1, "111", errors.New("222"), err1) //
//	log.Print("ERR2")
//	log.Printf("%ss:  %s", "%", err2)
//	log.Printf("%sv:  %v", "%", err2)
//	log.Printf("%s#v: %#v", "%", err2)
//
//	errors.New("222")
//}
