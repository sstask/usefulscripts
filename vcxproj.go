package vcxproj

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"github.com/axgle/mahonia"
)

var targetPath string
var dirList = make([]string, 0)
var fileListCpp = make([]string, 0)
var fileListH = make([]string, 0)
var fileListOther = make([]string, 0)

func isVaild(path string) bool {
	if path == "" || strings.Contains(path, "/.") || strings.Contains(path, "\\.") || len(path) == 0 {
		return false
	}
	return true
}

func getFilelist(path string) {
	err := filepath.Walk(path, func(path string, f os.FileInfo, err error) error {
		if f == nil {
			return err
		}
		if f.IsDir() && !isVaild(path) {
			return filepath.SkipDir
		}
		if f.IsDir() {
			dirList = append(dirList, path)
		} else {
			if strings.HasSuffix(path, ".cpp") {
				fileListCpp = append(fileListCpp, path)
			} else if strings.HasSuffix(path, ".h") {
				fileListH = append(fileListH, path)
			} else {
				fileListOther = append(fileListOther, path)
			}
		}
		return nil
	})
	if err != nil {
		fmt.Printf("filepath.Walk() returned %v\n", err)
	}
}

func createFilters(proname string) {
	fileName := proname + ".vcxproj.filters"
	fOut, err := os.Create(fileName)
	defer fOut.Close()
	if err != nil {
		fmt.Println("create file failed ", fileName, err)
	}
	strTemp := `<?xml version="1.0" encoding="utf-8"?>
<Project ToolsVersion="4.0" xmlns="http://schemas.microsoft.com/developer/msbuild/2003">`
	fOut.WriteString(strTemp)
	fOut.WriteString("\n <ItemGroup>")
	for _, v := range dirList {
		filtPath := strings.Replace(v, targetPath, "", -1)
		filtPath = strings.Replace(filtPath, "/", "\\", -1)
		if strings.HasPrefix(filtPath, "\\") {
			filtPath = filtPath[1:]
		}
		if filtPath == "" {
			continue
		}
		md5Code := md5.New()
		md5Code.Write([]byte(v))
		strTemp = "\n    <Filter Include=\"" + filtPath + "\">\n      <UniqueIdentifier>{" + hex.EncodeToString(md5Code.Sum(nil)) + "}</UniqueIdentifier>\n    </Filter>"
		fOut.WriteString(strTemp)
	}
	fOut.WriteString("\n </ItemGroup>")

	fOut.WriteString("\n <ItemGroup>")
	for _, v := range fileListCpp {
		fullPath := strings.Replace(v, "/", "\\", -1)
		filtPath := filepath.Dir(v)
		filtPath = strings.Replace(filtPath, targetPath, "", -1)
		filtPath = strings.Replace(filtPath, "/", "\\", -1)
		if strings.HasPrefix(filtPath, "\\") {
			filtPath = filtPath[1:]
		}
		strTemp = "\n    <ClCompile Include=\"" + fullPath + "\">\n      <Filter>" + filtPath + "</Filter>\n    </ClCompile>"
		fOut.WriteString(strTemp)
	}
	fOut.WriteString("\n </ItemGroup>")

	fOut.WriteString("\n <ItemGroup>")
	for _, v := range fileListH {
		fullPath := strings.Replace(v, "/", "\\", -1)
		filtPath := filepath.Dir(v)
		filtPath = strings.Replace(filtPath, targetPath, "", -1)
		filtPath = strings.Replace(filtPath, "/", "\\", -1)
		if strings.HasPrefix(filtPath, "\\") {
			filtPath = filtPath[1:]
		}
		strTemp = "\n    <ClInclude Include=\"" + fullPath + "\">\n      <Filter>" + filtPath + "</Filter>\n    </ClInclude>"
		fOut.WriteString(strTemp)
	}
	fOut.WriteString("\n </ItemGroup>")

	fOut.WriteString("\n <ItemGroup>")
	for _, v := range fileListOther {
		fullPath := strings.Replace(v, "/", "\\", -1)
		filtPath := filepath.Dir(v)
		filtPath = strings.Replace(filtPath, targetPath, "", -1)
		filtPath = strings.Replace(filtPath, "/", "\\", -1)
		if strings.HasPrefix(filtPath, "\\") {
			filtPath = filtPath[1:]
		}
		strTemp = "\n    <None Include=\"" + fullPath + "\">\n      <Filter>" + filtPath + "</Filter>\n    </None>"
		fOut.WriteString(strTemp)
	}
	fOut.WriteString("\n </ItemGroup>")

	fOut.WriteString("\n</Project>")
}

func createVcxproj(proname string) {
	fileName := proname + ".vcxproj"
	fOut, err := os.Create(fileName)
	defer fOut.Close()
	if err != nil {
		fmt.Println("create file failed ", fileName, err)
		return
	}
	md5Code := md5.New()
	md5Code.Write([]byte(targetPath))

	strTemp := `<?xml version="1.0" encoding="utf-8"?>
<Project DefaultTargets="Build" ToolsVersion="4.0" xmlns="http://schemas.microsoft.com/developer/msbuild/2003">
  <ItemGroup Label="ProjectConfigurations">
    <ProjectConfiguration Include="Debug|Win32">
      <Configuration>Debug</Configuration>
      <Platform>Win32</Platform>
    </ProjectConfiguration>
    <ProjectConfiguration Include="Release|Win32">
      <Configuration>Release</Configuration>
      <Platform>Win32</Platform>
    </ProjectConfiguration>
  </ItemGroup>
  <PropertyGroup Label="Globals">
    <ProjectGuid>{` + hex.EncodeToString(md5Code.Sum(nil)) + `}</ProjectGuid>
    <RootNamespace>hc</RootNamespace>
  </PropertyGroup>
  <Import Project="$(VCTargetsPath)\Microsoft.Cpp.Default.props" />
  <PropertyGroup Condition="'$(Configuration)|$(Platform)'=='Debug|Win32'" Label="Configuration">
    <ConfigurationType>Application</ConfigurationType>
    <UseDebugLibraries>true</UseDebugLibraries>
    <PlatformToolset>v110</PlatformToolset>
    <CharacterSet>MultiByte</CharacterSet>
  </PropertyGroup>
  <PropertyGroup Condition="'$(Configuration)|$(Platform)'=='Release|Win32'" Label="Configuration">
    <ConfigurationType>Application</ConfigurationType>
    <UseDebugLibraries>false</UseDebugLibraries>
    <PlatformToolset>v110</PlatformToolset>
    <WholeProgramOptimization>true</WholeProgramOptimization>
    <CharacterSet>MultiByte</CharacterSet>
  </PropertyGroup>
  <Import Project="$(VCTargetsPath)\Microsoft.Cpp.props" />
  <ImportGroup Label="ExtensionSettings">
  </ImportGroup>
  <ImportGroup Label="PropertySheets" Condition="'$(Configuration)|$(Platform)'=='Debug|Win32'">
    <Import Project="$(UserRootDir)\Microsoft.Cpp.$(Platform).user.props" Condition="exists('$(UserRootDir)\Microsoft.Cpp.$(Platform).user.props')" Label="LocalAppDataPlatform" />
  </ImportGroup>
  <ImportGroup Label="PropertySheets" Condition="'$(Configuration)|$(Platform)'=='Release|Win32'">
    <Import Project="$(UserRootDir)\Microsoft.Cpp.$(Platform).user.props" Condition="exists('$(UserRootDir)\Microsoft.Cpp.$(Platform).user.props')" Label="LocalAppDataPlatform" />
  </ImportGroup>
  <PropertyGroup Label="UserMacros" />
  <PropertyGroup />
  <ItemDefinitionGroup Condition="'$(Configuration)|$(Platform)'=='Debug|Win32'">
    <ClCompile>
      <WarningLevel>Level3</WarningLevel>
      <Optimization>Disabled</Optimization>
    </ClCompile>
    <Link>
      <GenerateDebugInformation>true</GenerateDebugInformation>
    </Link>
  </ItemDefinitionGroup>
  <ItemDefinitionGroup Condition="'$(Configuration)|$(Platform)'=='Release|Win32'">
    <ClCompile>
      <WarningLevel>Level3</WarningLevel>
      <Optimization>MaxSpeed</Optimization>
      <FunctionLevelLinking>true</FunctionLevelLinking>
      <IntrinsicFunctions>true</IntrinsicFunctions>
    </ClCompile>
    <Link>
      <GenerateDebugInformation>true</GenerateDebugInformation>
      <EnableCOMDATFolding>true</EnableCOMDATFolding>
      <OptimizeReferences>true</OptimizeReferences>
    </Link>
  </ItemDefinitionGroup>`
	fOut.WriteString(strTemp)

	fOut.WriteString("\n <ItemGroup>")
	for _, v := range fileListCpp {
		fullPath := strings.Replace(v, "/", "\\", -1)
		strTemp = "\n    <ClCompile Include=\"" + fullPath + "\" />"
		fOut.WriteString(strTemp)
	}
	fOut.WriteString("\n </ItemGroup>")

	fOut.WriteString("\n <ItemGroup>")
	for _, v := range fileListH {
		fullPath := strings.Replace(v, "/", "\\", -1)
		strTemp = "\n    <ClInclude Include=\"" + fullPath + "\" />"
		fOut.WriteString(strTemp)
	}
	fOut.WriteString("\n </ItemGroup>")

	fOut.WriteString("\n <ItemGroup>")
	for _, v := range fileListOther {
		fullPath := strings.Replace(v, "/", "\\", -1)
		strTemp = "\n    <None Include=\"" + fullPath + "\" />"
		fOut.WriteString(strTemp)
	}
	fOut.WriteString("\n </ItemGroup>")

	strTemp = `
  <Import Project="$(VCTargetsPath)\Microsoft.Cpp.targets" />
  <ImportGroup Label="ExtensionTargets">
  </ImportGroup>
</Project>`
	fOut.WriteString(strTemp)
}

func GenerateFile(vspath string, vsname string, dirpath string) {
	if dirpath == "." {
		dirpath, _ = filepath.Abs(filepath.Dir(os.Args[0]))
	}
	targetPath = dirpath
	getFilelist(dirpath)
	if !strings.HasSuffix(vspath, string(os.PathSeparator)) {
		vspath += string(os.PathSeparator) + vsname
	} else {
		vspath += vsname
	}
	createFilters(vspath)
	createVcxproj(vspath)
}

func VS() {
	dirpath, _ := filepath.Abs(filepath.Dir(os.Args[0]))
	cfgpath := dirpath + string(os.PathSeparator) + "cfg.json"
	fmt.Println(cfgpath)
	bytes, err := ioutil.ReadFile(cfgpath)
	if err != nil {
		fmt.Println("ReadFile: ", err.Error())
		return
	}

	decoder := mahonia.NewDecoder("gbk")
	var xxx = map[string]string{}
	if err := json.Unmarshal([]byte(decoder.ConvertString(string(bytes))), &xxx); err != nil {
		fmt.Println("Unmarshal: ", err.Error())
		return
	}
	GenerateFile(xxx["vspath"], xxx["vsname"], xxx["srcpath"])
}

//cfg.json
/*
{
"vspath":"D:\\svr\\svr",
"vsname":"svr",
"srcpath":"E:\\Server"
}
*/
