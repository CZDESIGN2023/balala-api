module go-cs

go 1.24.1

replace (
	cloud.google.com/go/compute/metadata => cloud.google.com/go/compute/metadata v0.2.3
	github.com/go-kratos/kratos/v2 => github.com/go-kratos/kratos/v2 v2.6.3
	google.golang.org/protobuf => google.golang.org/protobuf v1.31.0
)

require (
	github.com/99designs/gqlgen v0.17.22
	github.com/BurntSushi/toml v1.2.1
	github.com/Code-Hex/go-generics-cache v1.3.1
	github.com/PuerkitoBio/goquery v1.10.1
	github.com/RichardKnop/logging v0.0.0-20190827224416-1a693bdd4fae
	github.com/RichardKnop/machinery/v2 v2.0.11
	github.com/alibabacloud-go/darabonba-openapi/v2 v2.0.5
	github.com/alibabacloud-go/dysmsapi-20170525/v3 v3.0.6
	github.com/alibabacloud-go/tea v1.2.1
	github.com/alibabacloud-go/tea-utils/v2 v2.0.4
	github.com/aliyunmq/mq-http-go-sdk v1.0.3
	github.com/apache/pulsar-client-go v0.9.0
	github.com/apache/rocketmq-client-go/v2 v2.1.1
	github.com/apache/thrift v0.17.0
	github.com/casbin/casbin/v2 v2.77.2
	github.com/cloudwego/hertz v0.4.2
	github.com/dominikbraun/graph v0.23.0
	github.com/eclipse/paho.mqtt.golang v1.4.2
	github.com/eknkc/basex v1.0.1
	github.com/elastic/go-elasticsearch/v8 v8.13.1
	github.com/fanliao/go-promise v0.0.0-20141029170127-1890db352a72
	github.com/fasthttp/router v1.4.14
	github.com/gin-gonic/gin v1.9.1
	github.com/go-kratos/kratos/contrib/log/zap/v2 v2.0.0-20221019143716-dcae38d6568c
	github.com/go-kratos/kratos/v2 v2.7.0
	github.com/go-mysql-org/go-mysql v1.7.0
	github.com/go-redis/redis/v8 v8.11.5
	github.com/go-sql-driver/mysql v1.7.1
	github.com/go-stomp/stomp/v3 v3.0.5
	github.com/gogap/errors v0.0.0-20210818113853-edfbba0ddea9
	github.com/gogf/gf v1.16.9
	github.com/golang-jwt/jwt/v4 v4.5.0
	github.com/golang/mock v1.6.0
	github.com/gomodule/redigo v2.0.0+incompatible
	github.com/google/uuid v1.4.0
	github.com/google/wire v0.6.0
	github.com/gorilla/handlers v1.5.1
	github.com/gorilla/securecookie v1.1.2
	github.com/gorilla/websocket v1.5.0
	github.com/hibiken/asynq v0.24.1
	github.com/jhump/protoreflect v1.15.1
	github.com/joho/godotenv v1.5.1
	github.com/jordan-wright/email v4.0.1-0.20210109023952-943e75fe5223+incompatible
	github.com/lestrrat-go/file-rotatelogs v2.4.0+incompatible
	github.com/lucas-clemente/quic-go v0.31.1
	github.com/mitchellh/mapstructure v1.5.0
	github.com/mozillazg/go-pinyin v0.20.0
	github.com/nats-io/nats.go v1.22.1
	github.com/nicksnyder/go-i18n/v2 v2.2.1
	github.com/nsqio/go-nsq v1.1.0
	github.com/onsi/ginkgo v1.16.5
	github.com/onsi/gomega v1.24.1
	github.com/panjf2000/ants/v2 v2.8.2
	github.com/pingcap/errors v0.11.5-0.20240311024730-e056997136bb
	github.com/robfig/cron/v3 v3.0.1
	github.com/segmentio/kafka-go v0.4.38
	github.com/spf13/cast v1.3.1
	github.com/streadway/amqp v1.0.0
	github.com/stretchr/testify v1.9.0
	github.com/subosito/gotenv v1.2.0
	github.com/swaggo/swag v1.8.10
	github.com/tidwall/gjson v1.17.0
	github.com/tidwall/sjson v1.2.5
	github.com/u2takey/ffmpeg-go v0.5.0
	github.com/valyala/fasthttp v1.43.0
	github.com/vektah/gqlparser/v2 v2.5.1
	go.etcd.io/etcd/client/v3 v3.5.6
	go.mongodb.org/mongo-driver v1.4.6
	go.opentelemetry.io/otel/trace v1.24.0
	go.uber.org/zap v1.26.0
	golang.org/x/exp v0.0.0-20240205201215-2c58cdc269a3
	google.golang.org/genproto/googleapis/api v0.0.0-20240102182953-50ed04b92917
	google.golang.org/grpc v1.61.1
	google.golang.org/protobuf v1.34.0
	gopkg.in/vansante/go-ffprobe.v2 v2.2.1
	gorm.io/datatypes v1.2.0
	gorm.io/driver/mysql v1.5.2
	gorm.io/driver/postgres v1.5.0
	gorm.io/gorm v1.25.5
	gorm.io/plugin/dbresolver v1.4.7
	moul.io/zapgorm2 v1.1.3
)

require go.opentelemetry.io/otel/metric v1.24.0 // indirect

require (
	cloud.google.com/go v0.111.0 // indirect
	cloud.google.com/go/compute v1.23.3 // indirect
	cloud.google.com/go/compute/metadata v0.2.3 // indirect
	cloud.google.com/go/iam v1.1.5 // indirect
	cloud.google.com/go/pubsub v1.33.0 // indirect
	github.com/99designs/go-keychain v0.0.0-20191008050251-8e49817e8af4 // indirect
	github.com/99designs/keyring v1.2.1 // indirect
	github.com/AthenZ/athenz v1.10.39 // indirect
	github.com/DataDog/zstd v1.5.0 // indirect
	github.com/Knetic/govaluate v3.0.1-0.20171022003610-9aa49832a739+incompatible // indirect
	github.com/agnivade/levenshtein v1.1.1 // indirect
	github.com/alibabacloud-go/alibabacloud-gateway-spi v0.0.4 // indirect
	github.com/alibabacloud-go/debug v0.0.0-20190504072949-9472017b5c68 // indirect
	github.com/alibabacloud-go/endpoint-util v1.1.0 // indirect
	github.com/alibabacloud-go/openapi-util v0.1.0 // indirect
	github.com/alibabacloud-go/tea-utils v1.3.1 // indirect
	github.com/alibabacloud-go/tea-xml v1.1.3 // indirect
	github.com/aliyun/credentials-go v1.3.1 // indirect
	github.com/andybalholm/brotli v1.0.4 // indirect
	github.com/andybalholm/cascadia v1.3.3 // indirect
	github.com/ardielle/ardielle-go v1.5.2 // indirect
	github.com/aws/aws-sdk-go v1.55.5 // indirect
	github.com/beorn7/perks v1.0.1 // indirect
	github.com/bufbuild/protocompile v0.4.0 // indirect
	github.com/bytedance/go-tagexpr/v2 v2.9.2 // indirect
	github.com/bytedance/gopkg v0.0.0-20220413063733-65bf48ffb3a7 // indirect
	github.com/bytedance/sonic v1.11.6 // indirect
	github.com/bytedance/sonic/loader v0.1.1 // indirect
	github.com/clbanning/mxj/v2 v2.5.5 // indirect
	github.com/cloudwego/base64x v0.1.4 // indirect
	github.com/cloudwego/iasm v0.2.0 // indirect
	github.com/cloudwego/netpoll v0.3.1 // indirect
	github.com/danieljoos/wincred v1.1.2 // indirect
	github.com/dvsekhvalnov/jose2go v1.5.0 // indirect
	github.com/elastic/elastic-transport-go/v8 v8.5.0 // indirect
	github.com/emirpasic/gods v1.12.0 // indirect
	github.com/facebookgo/stack v0.0.0-20160209184415-751773369052 // indirect
	github.com/felixge/httpsnoop v1.0.1 // indirect
	github.com/gabriel-vasile/mimetype v1.4.3 // indirect
	github.com/go-kratos/aegis v0.2.0 // indirect
	github.com/go-redsync/redsync/v4 v4.0.4 // indirect
	github.com/go-stack/stack v1.8.0 // indirect
	github.com/go-task/slim-sprig v0.0.0-20210107165309-348f09dbbbc0 // indirect
	github.com/godbus/dbus v0.0.0-20190726142602-4481cbc300e2 // indirect
	github.com/gogap/stack v0.0.0-20150131034635-fef68dddd4f8 // indirect
	github.com/golang-jwt/jwt v3.2.1+incompatible // indirect
	github.com/golang/groupcache v0.0.0-20210331224755-41bb18bfe9da // indirect
	github.com/golang/snappy v0.0.4 // indirect
	github.com/google/pprof v0.0.0-20210720184732-4bb14d4b1be1 // indirect
	github.com/google/s2a-go v0.1.7 // indirect
	github.com/googleapis/enterprise-certificate-proxy v0.3.2 // indirect
	github.com/googleapis/gax-go/v2 v2.12.0 // indirect
	github.com/gsterjov/go-libsecret v0.0.0-20161001094733-a6f4afe4910c // indirect
	github.com/hashicorp/errwrap v1.1.0 // indirect
	github.com/hashicorp/go-multierror v1.1.1 // indirect
	github.com/hashicorp/golang-lru v1.0.2 // indirect
	github.com/henrylee2cn/ameda v1.4.10 // indirect
	github.com/henrylee2cn/goutil v0.0.0-20210127050712-89660552f6f8 // indirect
	github.com/jackc/pgx/v5 v5.3.0 // indirect
	github.com/jmespath/go-jmespath v0.4.0 // indirect
	github.com/kelseyhightower/envconfig v1.4.0 // indirect
	github.com/klauspost/compress v1.17.1 // indirect
	github.com/klauspost/cpuid/v2 v2.2.7 // indirect
	github.com/linkedin/goavro/v2 v2.9.8 // indirect
	github.com/marten-seemann/qpack v0.3.0 // indirect
	github.com/marten-seemann/qtls-go1-18 v0.1.3 // indirect
	github.com/marten-seemann/qtls-go1-19 v0.1.1 // indirect
	github.com/mattn/go-sqlite3 v1.14.16 // indirect
	github.com/matttproud/golang_protobuf_extensions v1.0.1 // indirect
	github.com/mtibben/percent v0.2.1 // indirect
	github.com/nats-io/nats-server/v2 v2.9.10 // indirect
	github.com/nats-io/nkeys v0.3.0 // indirect
	github.com/nats-io/nuid v1.0.1 // indirect
	github.com/nyaruka/phonenumbers v1.0.55 // indirect
	github.com/opentracing/opentracing-go v1.2.0 // indirect
	github.com/patrickmn/go-cache v2.1.0+incompatible // indirect
	github.com/pierrec/lz4 v2.0.5+incompatible // indirect
	github.com/pierrec/lz4/v4 v4.1.17 // indirect
	github.com/pingcap/log v1.1.0 // indirect
	github.com/pingcap/tidb/parser v0.0.0-20221126021158-6b02a5d8ba7d // indirect
	github.com/prometheus/client_golang v1.11.1 // indirect
	github.com/prometheus/client_model v0.4.0 // indirect
	github.com/prometheus/common v0.26.0 // indirect
	github.com/prometheus/procfs v0.6.0 // indirect
	github.com/redis/go-redis/v9 v9.0.3 // indirect
	github.com/satori/go.uuid v1.2.0 // indirect
	github.com/savsgio/gotils v0.0.0-20220530130905-52f3993e8d6d // indirect
	github.com/shopspring/decimal v0.0.0-20180709203117-cd690d0c9e24 // indirect
	github.com/siddontang/go v0.0.0-20180604090527-bdc77568d726 // indirect
	github.com/siddontang/go-log v0.0.0-20180807004314-8d05993dda07 // indirect
	github.com/sirupsen/logrus v1.9.3 // indirect
	github.com/spaolacci/murmur3 v1.1.0 // indirect
	github.com/stathat/consistent v1.0.0 // indirect
	github.com/tidwall/match v1.1.1 // indirect
	github.com/tidwall/pretty v1.2.1 // indirect
	github.com/tjfoc/gmsm v1.3.2 // indirect
	github.com/twitchyliquid64/golang-asm v0.15.1 // indirect
	github.com/u2takey/go-utils v0.3.1 // indirect
	github.com/valyala/bytebufferpool v1.0.0 // indirect
	github.com/xdg/scram v1.0.5 // indirect
	github.com/xdg/stringprep v1.0.3 // indirect
	go.opencensus.io v0.24.0 // indirect
	go.uber.org/goleak v1.3.0 // indirect
	golang.org/x/arch v0.7.0 // indirect
	golang.org/x/mod v0.17.0 // indirect
	golang.org/x/oauth2 v0.15.0 // indirect
	golang.org/x/term v0.27.0 // indirect
	golang.org/x/time v0.3.0 // indirect
	google.golang.org/api v0.149.0 // indirect
	google.golang.org/appengine v1.6.8 // indirect
	google.golang.org/genproto v0.0.0-20231212172506-995d672761c0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20240102182953-50ed04b92917 // indirect
	gopkg.in/ini.v1 v1.62.0 // indirect
	gopkg.in/natefinch/lumberjack.v2 v2.2.1 // indirect
	gorm.io/driver/sqlite v1.5.0 // indirect
)

require (
	github.com/KyleBanks/depth v1.2.1 // indirect
	github.com/cespare/xxhash/v2 v2.2.0
	github.com/coreos/go-semver v0.3.0 // indirect
	github.com/coreos/go-systemd/v22 v22.3.2 // indirect
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/dgryski/go-rendezvous v0.0.0-20200823014737-9f7001d12a5f // indirect
	github.com/fsnotify/fsnotify v1.6.0 // indirect
	github.com/gin-contrib/sse v0.1.0 // indirect
	github.com/go-kratos/kratos/contrib/registry/etcd/v2 v2.0.0-20230410063538-3958f9d5c06e
	github.com/go-logr/logr v1.4.1 // indirect
	github.com/go-logr/stdr v1.2.2 // indirect
	github.com/go-openapi/jsonpointer v0.19.6 // indirect
	github.com/go-openapi/jsonreference v0.20.2 // indirect
	github.com/go-openapi/spec v0.20.8 // indirect
	github.com/go-openapi/swag v0.22.3 // indirect
	github.com/go-playground/form/v4 v4.2.1 // indirect
	github.com/go-playground/locales v0.14.1 // indirect
	github.com/go-playground/universal-translator v0.18.1 // indirect
	github.com/go-playground/validator/v10 v10.20.0
	github.com/goccy/go-json v0.10.2 // indirect
	github.com/gogo/protobuf v1.3.2 // indirect
	github.com/golang/protobuf v1.5.3
	github.com/google/go-cmp v0.6.0
	github.com/gorilla/mux v1.8.0
	github.com/imdario/mergo v0.3.13 // indirect
	github.com/jackc/pgpassfile v1.0.0 // indirect
	github.com/jackc/pgservicefile v0.0.0-20221227161230-091c0ba34f0a // indirect
	github.com/jinzhu/inflection v1.0.0 // indirect
	github.com/jinzhu/now v1.1.5 // indirect
	github.com/jonboulle/clockwork v0.3.0 // indirect
	github.com/josharian/intern v1.0.0 // indirect
	github.com/json-iterator/go v1.1.12 // indirect
	github.com/leodido/go-urn v1.4.0 // indirect
	github.com/lestrrat-go/strftime v1.0.6 // indirect
	github.com/mailru/easyjson v0.7.7 // indirect
	github.com/mattn/go-isatty v0.0.20 // indirect
	github.com/modern-go/concurrent v0.0.0-20180306012644-bacd9c7ef1dd // indirect
	github.com/modern-go/reflect2 v1.0.2 // indirect
	github.com/nxadm/tail v1.4.8 // indirect
	github.com/onsi/ginkgo/v2 v2.6.1 // indirect
	github.com/openzipkin/zipkin-go v0.4.1 // indirect
	github.com/pelletier/go-toml/v2 v2.2.1 // indirect
	github.com/pkg/errors v0.9.1 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	github.com/ugorji/go/codec v1.2.12 // indirect
	go.etcd.io/etcd/api/v3 v3.5.6 // indirect
	go.etcd.io/etcd/client/pkg/v3 v3.5.6 // indirect
	go.opentelemetry.io/otel v1.24.0
	go.opentelemetry.io/otel/exporters/jaeger v1.11.2
	go.opentelemetry.io/otel/exporters/zipkin v1.11.2
	go.opentelemetry.io/otel/sdk v1.24.0
	go.uber.org/atomic v1.11.0 // indirect
	go.uber.org/multierr v1.11.0 // indirect
	golang.org/x/crypto v0.31.0 // indirect
	golang.org/x/net v0.33.0 // indirect
	golang.org/x/sync v0.10.0 // indirect
	golang.org/x/sys v0.28.0 // indirect
	golang.org/x/text v0.21.0
	golang.org/x/tools v0.21.1-0.20240508182429-e35e4ccd0d2d // indirect
	gopkg.in/tomb.v1 v1.0.0-20141024135613-dd632973f1e7 // indirect
	gopkg.in/yaml.v2 v2.4.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)
