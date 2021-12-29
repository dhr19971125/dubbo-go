package common

import (
	"dubbo.apache.org/dubbo-go/v3/common/constant"
	"fmt"
	"hash/crc32"
	"sort"
	"testing"
)

import (
	"github.com/stretchr/testify/assert"
)

func TestMetadataInfoCalAndGetRevision(t *testing.T)  {

	serviceUrl, _ := NewURL("dubbo://127.0.0.1:20000/com.ikurento.user.UserProvider?anyhost=true&" +
		"application=BDTService&category=providers&default.timeout=10000&dubbo=dubbo-provider-golang-1.0.0&" +
		"environment=dev&interface=com.ikurento.user.UserProvider&ip=192.168.56.1&methods=GetUser%2C&" +
		"module=dubbogo+user-info+server&org=ikurento.com&owner=ZX&pid=1447&revision=0.0.1&" +
		"side=provider&timeout=3000&timestamp=155650979798")
	metadataInfo := NewMetadataInfo("appTest","revisionTest", map[string]*ServiceInfo{
		"com.ikurento.user.UserProvider": NewServiceInfoWithURL(serviceUrl),
	})

	t.Run("RevisionAndReportedExist", func(t *testing.T) {
		metadataInfo.MarkReported()
		expected:="revisionTest"
		res:=metadataInfo.CalAndGetRevision()
		assert.Equal(t,expected,res)
	})

	t.Run("ServicesBlank", func(t *testing.T) {
		if metadataInfo.HasReported()!=false{
			metadataInfo.Reported=false
		}
		expected:="0"
		metadataInfo.Services=nil
		res:=metadataInfo.CalAndGetRevision()
		assert.Equal(t,res,expected)
	})

	t.Run("msBlank", func(t *testing.T) {
		candidates := make([]string, 8)
		for _, s := range metadataInfo.Services {
			sk := s.ServiceKey
			candidates = append(candidates, sk)
		}
		sort.Strings(candidates)
		expected := uint64(0)
		for _, c := range candidates {
			expected += uint64(crc32.ChecksumIEEE([]byte(c)))
		}
		if metadataInfo.HasReported()!=false{
			metadataInfo.Reported=false
		}
		res:=metadataInfo.CalAndGetRevision()
		assert.Equal(t,fmt.Sprint(expected),res)
	})

	t.Run("msExist", func(t *testing.T) {

		for _, s := range metadataInfo.Services {
			s.URL.Methods=[]string{"test"}
		}
		res:=metadataInfo.CalAndGetRevision()

		candidates := make([]string, 8)
		for _, s := range metadataInfo.Services {
			sk := s.ServiceKey
			ms := s.URL.Methods

				for _, m := range ms {
					// methods are part of candidates
					candidates = append(candidates, sk+constant.KeySeparator+m)
				}
			}
		sort.Strings(candidates)
		expected := uint64(0)
		for _, c := range candidates {
			expected += uint64(crc32.ChecksumIEEE([]byte(c)))
		}
		assert.Equal(t,fmt.Sprint(expected),res)
	})
}

func TestMetadataInfoAddService(t *testing.T)  {

	t.Run("serviceNil", func(t *testing.T) {
		metadataInfo := NewMetadataInfWithApp("appTest")
		metadataInfo.AddService(metadataInfo.Services["Services"])
		assert.Nil(t,metadataInfo.Services["Services"])
	})

	t.Run("AddService", func(t *testing.T) {
		serviceUrl, _ := NewURL("dubbo://127.0.0.1:20000/com.ikurento.user.UserProvider?anyhost=true&" +
			"application=BDTService&category=providers&default.timeout=10000&dubbo=dubbo-provider-golang-1.0.0&" +
			"environment=dev&interface=com.ikurento.user.UserProvider&ip=192.168.56.1&methods=GetUser%2C&" +
			"module=dubbogo+user-info+server&org=ikurento.com&owner=ZX&pid=1447&revision=0.0.1&" +
			"side=provider&timeout=3000&timestamp=155650979798")
		serviceTest:=NewServiceInfoWithURL(serviceUrl)
		metadataInfo := NewMetadataInfWithApp("appTest")
		metadataInfo.AddService(serviceTest)
		assert.Equal(t,serviceTest,metadataInfo.Services[serviceTest.GetMatchKey()])
	})
}

func TestMetadataInfoRemoveService(t *testing.T)  {
	t.Run("serviceNil", func(t *testing.T) {
		metadataInfo := NewMetadataInfWithApp("appTest")
		metadataInfo.RemoveService(metadataInfo.Services["Services"])
		assert.Nil(t,metadataInfo.Services["Services"])

	})

	t.Run("RemoveService", func(t *testing.T) {
		serviceUrl, _ := NewURL("dubbo://127.0.0.1:20000/com.ikurento.user.UserProvider?anyhost=true&" +
			"application=BDTService&category=providers&default.timeout=10000&dubbo=dubbo-provider-golang-1.0.0&" +
			"environment=dev&interface=com.ikurento.user.UserProvider&ip=192.168.56.1&methods=GetUser%2C&" +
			"module=dubbogo+user-info+server&org=ikurento.com&owner=ZX&pid=1447&revision=0.0.1&" +
			"side=provider&timeout=3000&timestamp=155650979798")
		serviceInfoTest:= NewServiceInfoWithURL(serviceUrl)
		metadataInfo := NewMetadataInfo("appTest", "revisionTest", map[string]*ServiceInfo{
			serviceInfoTest.GetMatchKey():serviceInfoTest})
		metadataInfo.RemoveService(serviceInfoTest)
		assert.Nil(t,metadataInfo.Services[serviceInfoTest.GetMatchKey()])
	})
}

func TestServiceInfoGetMethods(t *testing.T)  {
	t.Run("MethodsKeyBlank", func(t *testing.T) {
		url, _ := NewURL("dubbo://127.0.0.1:20000/com.ikurento.user.UserProvider?anyhost=true&" +
			"application=BDTService&category=providers&default.timeout=10000&dubbo=dubbo-provider-golang-1.0.0&" +
			"environment=dev&interface=com.ikurento.user.UserProvider&ip=192.168.56.1&methods=GetUser%2C&" +
			"module=dubbogo+user-info+server&org=ikurento.com&owner=ZX&pid=1447&revision=0.0.1&" +
			"side=provider&timeout=3000&timestamp=155650979798")
		val:=url.GetParam("methods","")
		serviceInfoTest:= NewServiceInfo(url.Service(),url.Group(), url.Version(), url.Protocol, url.Path, map[string]string{"methods":val})
		//params:=serviceInfoTest.GetParams()
		//res:=params["methods"]
		Params:=serviceInfoTest.GetParams()
		expected:=Params["methods"][0]
		res:=serviceInfoTest.GetMethods()
		assert.Equal(t,expected[:len(expected)-1],res[0])

	})

	t.Run("GetMethods", func(t *testing.T) {
		serviceUrl, _ := NewURL("dubbo://127.0.0.1:20000/com.ikurento.user.UserProvider?anyhost=true&" +
			"application=BDTService&category=providers&default.timeout=10000&dubbo=dubbo-provider-golang-1.0.0&" +
			"environment=dev&interface=com.ikurento.user.UserProvider&ip=192.168.56.1&methods=GetUser%2C&" +
			"module=dubbogo+user-info+server&org=ikurento.com&owner=ZX&pid=1447&revision=0.0.1&" +
			"side=provider&timeout=3000&timestamp=155650979798")

		serviceInfoTest:= NewServiceInfoWithURL(serviceUrl)
		expect:=serviceInfoTest.MatchKey
		serviceInfoTest.MatchKey=""
		MatchKey:=serviceInfoTest.GetMatchKey()
		assert.Equal(t,expect,MatchKey)
	})
}

func TestServiceInfoGetServiceKey(t *testing.T)  {
	t.Run("ServiceKeyNotBlank", func(t *testing.T) {
		serviceUrl, _ := NewURL("dubbo://127.0.0.1:20000/com.ikurento.user.UserProvider?anyhost=true&" +
			"application=BDTService&category=providers&default.timeout=10000&dubbo=dubbo-provider-golang-1.0.0&" +
			"environment=dev&interface=com.ikurento.user.UserProvider&ip=192.168.56.1&methods=GetUser%2C&" +
			"module=dubbogo+user-info+server&org=ikurento.com&owner=ZX&pid=1447&revision=0.0.1&" +
			"side=provider&timeout=3000&timestamp=155650979798")
		serviceInfoTest:= NewServiceInfoWithURL(serviceUrl)
		serviceKey:=serviceInfoTest.GetServiceKey()
		assert.Equal(t,serviceInfoTest.ServiceKey,serviceKey)
	})

	t.Run("ServiceKeyBlank", func(t *testing.T) {
		serviceUrl, _ := NewURL("dubbo://127.0.0.1:20000/com.ikurento.user.UserProvider?anyhost=true&" +
			"application=BDTService&category=providers&default.timeout=10000&dubbo=dubbo-provider-golang-1.0.0&" +
			"environment=dev&interface=com.ikurento.user.UserProvider&ip=192.168.56.1&methods=GetUser%2C&" +
			"module=dubbogo+user-info+server&org=ikurento.com&owner=ZX&pid=1447&revision=0.0.1&" +
			"side=provider&timeout=3000&timestamp=155650979798")

		serviceInfoTest:= NewServiceInfoWithURL(serviceUrl)
		expect:=serviceInfoTest.ServiceKey
		serviceInfoTest.ServiceKey=""
		serviceKey:=serviceInfoTest.GetServiceKey()
		assert.Equal(t,expect,serviceKey)
	})
}
