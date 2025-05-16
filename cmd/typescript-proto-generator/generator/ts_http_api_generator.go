package generator

import (
	parser2 "go-cs/pkg/parser"
	"os"
	"path/filepath"
	"text/template"
)

const httpApiDartTemplate = `// Auto-generated Typescript code. DO NOT MODIFY.
// Source: {{.InputFilePath}}
import 'package:flutter_net/http.dart';
import 'package:grpc/grpc.dart';

import '../data_center.dart';
import '../pb/bean/request.pb.dart';
import '{{GetImportGrpcPathByPorto .InputFilePath}}';
import '../pb/comm/errors.pb.dart' as comm;
import '../utils/error.dart';
{{$sericeName := GetTitleNameByPath .InputFilePath}}
class {{ $sericeName }}Api {
  {{range .Services}}{{ $serviceName := .Name }}
  {{range .Methods}}{{ $methodName := .Name }}{{ $inputType := .InputType }}{{ $outputType := .OutputType }}{{ $PostPath0 := index .PostPath 0 }}
  {{- if eq .OutputType "CommonReply"}}{{ $outputType = "comm.CommonReply" }}{{- end }}
  static Future<{{ $outputType }}> {{ ToHump $methodName }}({{ $inputType }} data,
    {bool? socketHttp, bool? http2, String method = "{{ $PostPath0.Method }}"}) async {
    // 访问路径
    var paths = { {{ range .PostPath }}
      "{{ .Method }}": "{{ GetDartApiPath .FormatString }}",
   {{end}} }; 
    var path = paths[method];
    if (path == null) {
      return _on{{ $methodName }}Error("{{ $serviceName }}.{{ ToHump $methodName }} method:$method Not defined!");
    }
    // socket模拟短连接
    if (socketHttp ?? dataCenter.useSocketHttp) {
      return _{{ ToHump $methodName }}Socket(path, method, data);
    }
    // http2 - grpc
    if (http2 ?? dataCenter.useHttp2) {
      return _{{ ToHump $methodName }}Grpc(data);
    }
    // 普通http
    late {{ $outputType }} response;
    try {
      response = await HttpRequest.sendProbuf<{{ $outputType }}>(
          dataCenter.getHttpUrl(path), 
          data, 
          {{ $outputType }}(), 
          ProtoBufError(), 
          method: method);
      dataCenter.apiDebugPrint('Greeter client received: $response');
    } catch (e) {
      response = _on{{ $methodName }}Error(e);
    }
    return response;
  }
  // http2 - grpc
  static Future<{{ $outputType }}> _{{ ToHump $methodName }}Grpc({{ $inputType }} data) async {
    final channel = ClientChannel(
      dataCenter.serverHost,
      port: dataCenter.serverPortRpc,
      options: const ChannelOptions(credentials: ChannelCredentials.insecure()),
    );
    final stub = {{ $serviceName }}Client(channel);
    late {{ $outputType }} response;
    try {
      response = await stub.{{ ToHump $methodName }}(data);
      dataCenter.apiDebugPrint('Greeter client received: $response');
    } catch (e) {
      response = _on{{ $methodName }}Error(e);
    }
    await channel.shutdown();
    return response;
  }
  // socket模拟短连接
  static Future<{{ $outputType }}> _{{ ToHump $methodName }}Socket(String method, String path, {{ $inputType }} data) async {
    late {{ $outputType }} response;
    try {
      Request request = Request();
      request.path = path;
      switch (method) {
        case HttpRequest.methodGet:
          request.method = Request_Method.GET;
          break;
        case HttpRequest.methodPost:
          request.method = Request_Method.POST;
          request.data = data.writeToBuffer();
          break;
        case HttpRequest.methodPut:
          request.method = Request_Method.PUT;
          break;
        case HttpRequest.methodDelete:
          request.method = Request_Method.DELETE;
          break;
      }
      response = {{ $outputType }}();
      await dataCenter.socketMgr.httpTunnel.protoBufHttpRequest(request, response);
    } catch (e) {
      response = _on{{ $methodName }}Error(e);
    }
    return response;
  }
  static {{ $outputType }} _on{{ $methodName }}Error(e) {
    dataCenter.apiDebugPrint('Caught error: $e');
    int? code = StatusCode.unknown;
    String? message;
    if (e is GrpcError) {
      code = e.code;
      message = e.message;
    } else if (e is String) {
      message = e;
    } else if (e is Exception) {
      message = e.toString();
    } else if (e is Error) {
      message = e.toString();
    }
    var response = {{ $outputType }}()
      ..error = (comm.ErrorInfo()
        ..code = code
        ..message = message ?? '');
    return response;
  }
  {{end}}{{end}}
}
`

func GenerateDartHttpApiFile(tmplParams *parser2.TemplateParams, outputFilePath string) error {
	funcMap := template.FuncMap{
		"ToUpperWithUnderscores":   parser2.ToUpperWithUnderscores,
		"ToPascalCase":             parser2.ToPascalCase,
		"ToHump":                   parser2.ToHump,
		"GetNameByPath":            parser2.GetNameByPath,
		"GetTitleNameByPath":       parser2.GetTitleNameByPath,
		"GetImportGrpcPathByPorto": parser2.GetImportGrpcPathByPorto,
		"GetDartApiPath":           parser2.GetDartApiPath,
	}

	tmpl, err := template.New("httpApiDart").Funcs(funcMap).Parse(httpApiDartTemplate)
	if err != nil {
		return err
	}
	// 计算 EnumTypePrefix
	tmplParams.EnumTypePrefix = parser2.ComputeEnumTypePrefix(tmplParams.Messages)

	// 创建输出文件所在的目录
	dir := filepath.Dir(outputFilePath)
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		err = os.MkdirAll(dir, 0755) // Creates directory with permissions set to 0755
		if err != nil {
			return err
		}
	}

	outputFile, err := os.Create(outputFilePath)
	if err != nil {
		return err
	}
	defer outputFile.Close()

	err = tmpl.Execute(outputFile, tmplParams)
	if err != nil {
		return err
	}

	return nil
}
