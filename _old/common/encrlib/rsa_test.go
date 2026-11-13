package encrlib

//func TestNewRSAPrivateKey(t *testing.T) {
//
//	pathToStore := "test_rsa_key_" + strconv.FormatInt(time.Now().UnixNano(), 10) + ".test"
//	privateKey, err := NewRSAPrivateKey(pathToStore)
//	require.NoError(t, err)
//	require.NotNil(t, privateKey)
//
//	privateKeyAgain, err := NewRSAPrivateKey(pathToStore)
//	require.NoError(t, err)
//	require.NotNil(t, privateKeyAgain)
//
//	//require.Equal(t, *privateKey, *privateKeyAgain)
//	require.True(t, reflect.DeepEqual(*privateKey, *privateKeyAgain))
//
//}
