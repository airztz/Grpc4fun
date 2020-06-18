#curl -k -v "http://localhost:8081/hello" -d '{"featureName": "Joy", "complexFeatureValue": {
#  "@type": "google.protobuf.Int32Value",
#  "value": 111
#}}'
#curl -k -v "http://localhost:8081/hello" -d '{"featureName": "Joy", "complexFeatureValue": {
#  "@type": "google.protobuf.FloatValue",
#  "value": 1.2345
#}}'
#curl -k -v "http://localhost:8081/hello" -d '{"featureName": "Joy", "complexFeatureValue": {
#  "@type": "google.protobuf.ListValue",
#  "value": [1,2.2,"3", false]
#}}'
#curl -k -v "http://localhost:8081/hello" -d '{"featureName": "Joy", "complexFeatureValue": {
#  "@type": "google.protobuf.ListValue",
#  "value": [1,2.2,"3",[1,2,3], {"a":"a", "b":1}]
#}}'
curl -k -v "http://localhost:8081/hello" -d '{"featureName": "Joy", "complexFeatureValue": {
  "@type": "google.protobuf.Struct",
  "value": {"foo": 123, "bar": "bar", "baz": 1.23, "zoo": [1,2.2,"3",false]}
}}'
#curl -k -v "http://localhost:8081/hello" -d '{"featureName": "Joy", "complexFeatureValue": {
#  "@type": "google.protobuf.Struct",
#  "value": {"foo": 123, "bar": "bar", "baz": {"a":"a", "b":1}}
#}}'

#curl -k -v "http://localhost:8081/hello" -d '{"featureName": "Joy", "complexStructValue": {"foo": 123, "bar": "bar", "baz": {"a":"a", "b":1}}}'