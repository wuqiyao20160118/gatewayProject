(window["webpackJsonp"]=window["webpackJsonp"]||[]).push([["chunk-f860e2e0"],{"22ce":function(e,t,i){"use strict";i.d(t,"f",(function(){return a})),i.d(t,"d",(function(){return l})),i.d(t,"b",(function(){return s})),i.d(t,"i",(function(){return o})),i.d(t,"c",(function(){return c})),i.d(t,"j",(function(){return n})),i.d(t,"a",(function(){return u})),i.d(t,"h",(function(){return m})),i.d(t,"e",(function(){return d})),i.d(t,"g",(function(){return f}));var r=i("b775");function a(e){return Object(r["a"])({url:"/service/service_list",method:"get",params:e})}function l(e){return Object(r["a"])({url:"/service/service_delete",method:"get",params:e})}function s(e){return Object(r["a"])({url:"/service/service_add_http",method:"post",data:e})}function o(e){return Object(r["a"])({url:"/service/service_update_http",method:"post",data:e})}function c(e){return Object(r["a"])({url:"/service/service_add_tcp",method:"post",data:e})}function n(e){return Object(r["a"])({url:"/service/service_update_tcp",method:"post",data:e})}function u(e){return Object(r["a"])({url:"/service/service_add_grpc",method:"post",data:e})}function m(e){return Object(r["a"])({url:"/service/service_update_grpc",method:"post",data:e})}function d(e){return Object(r["a"])({url:"/service/service_detail",method:"get",params:e})}function f(e){return Object(r["a"])({url:"/service/service_stat",method:"get",params:e})}},"37bc":function(e,t,i){"use strict";var r=i("a4d4"),a=i.n(r);a.a},8446:function(e,t,i){"use strict";i.r(t);var r=function(){var e=this,t=e.$createElement,i=e._self._c||t;return i("div",{staticClass:"mixin-components-container"},[i("el-row",[i("el-card",{staticClass:"box-card"},[i("div",{staticClass:"clearfix",attrs:{slot:"header"},slot:"header"},[!1===e.isEdit?i("span",[e._v("Create TCP Service")]):e._e(),!0===e.isEdit?i("span",[e._v("Edit TCP Service")]):e._e()]),i("div",{staticStyle:{"margin-bottom":"50px"}},[i("el-form",{ref:"form",attrs:{model:e.form,"label-width":"140px"}},[i("el-form-item",{staticClass:"is-required",attrs:{label:"Service Name"}},[i("el-input",{attrs:{placeholder:"6-128 characters",disabled:!0===e.isEdit},model:{value:e.form.service_name,callback:function(t){e.$set(e.form,"service_name",t)},expression:"form.service_name"}})],1),i("el-form-item",{staticClass:"is-required",attrs:{label:"Service Discription"}},[i("el-input",{attrs:{placeholder:"0-255 characters, required"},model:{value:e.form.service_desc,callback:function(t){e.$set(e.form,"service_desc",t)},expression:"form.service_desc"}})],1),i("el-form-item",{staticClass:"is-required",attrs:{label:"Port"}},[i("el-input",{attrs:{placeholder:"The input should be an integer between 8001-8999, required",disabled:!0===e.isEdit},model:{value:e.form.port,callback:function(t){e.$set(e.form,"port",t)},expression:"form.port"}})],1),i("el-form-item",{attrs:{label:"Enable Validation"}},[i("el-switch",{attrs:{"active-value":1,"inactive-value":0},model:{value:e.form.open_auth,callback:function(t){e.$set(e.form,"open_auth",t)},expression:"form.open_auth"}})],1),i("el-form-item",{attrs:{label:"IP Whitelist"}},[i("el-input",{attrs:{type:"textarea",placeholder:"Format: 127.0.0.1:80, multiple lines, whitelist is prior to blacklist",autosize:""},model:{value:e.form.white_list,callback:function(t){e.$set(e.form,"white_list",t)},expression:"form.white_list"}})],1),i("el-form-item",{attrs:{label:"IP Blacklist"}},[i("el-input",{attrs:{type:"textarea",placeholder:"Multiple lines",autosize:""},model:{value:e.form.black_list,callback:function(t){e.$set(e.form,"black_list",t)},expression:"form.black_list"}})],1),i("el-form-item",{attrs:{label:"Flow Limitation on ClientIP"}},[i("el-input",{attrs:{placeholder:"0: no limitation"},model:{value:e.form.clientip_flow_limit,callback:function(t){e.$set(e.form,"clientip_flow_limit",t)},expression:"form.clientip_flow_limit"}})],1),i("el-form-item",{attrs:{label:"Flow Limitation on Service"}},[i("el-input",{attrs:{placeholder:"0: no limitation"},model:{value:e.form.service_flow_limit,callback:function(t){e.$set(e.form,"service_flow_limit",t)},expression:"form.service_flow_limit"}})],1),i("el-form-item",{attrs:{label:"Round Type"}},[i("el-radio-group",{model:{value:e.form.round_type,callback:function(t){e.$set(e.form,"round_type",t)},expression:"form.round_type"}},[i("el-radio",{attrs:{label:"0"}},[e._v("random")]),i("el-radio",{attrs:{label:"1"}},[e._v("round-robin")]),i("el-radio",{attrs:{label:"2"}},[e._v("weight_round-robin")]),i("el-radio",{attrs:{label:"3"}},[e._v("ip_hash")])],1)],1),i("el-form-item",{staticClass:"is-required",attrs:{label:"IP List"}},[i("el-input",{attrs:{type:"textarea",placeholder:"Format: 127.0.0.1:80, multiple lines",autosize:""},model:{value:e.form.ip_list,callback:function(t){e.$set(e.form,"ip_list",t)},expression:"form.ip_list"}})],1),i("el-form-item",{staticClass:"is-required",attrs:{label:"Weight List"}},[i("el-input",{attrs:{type:"textarea",placeholder:"Multiple lines",autosize:""},model:{value:e.form.weight_list,callback:function(t){e.$set(e.form,"weight_list",t)},expression:"form.weight_list"}})],1),i("el-form-item",[i("el-button",{attrs:{type:"primary",disabled:e.submitButtonDisabled},on:{click:e.onSubmit}},[e._v("Submit")])],1)],1)],1)])],1)],1)},a=[],l=(i("a9e3"),i("d3b7"),i("ac1f"),i("25f0"),i("5319"),i("22ce")),s={name:"ServiceAddTcp",data:function(){return{isEdit:!1,submitButtonDisabled:!1,form:{id:"0",service_name:"",service_desc:"",port:"",round_type:2,ip_list:"",weight_list:"",open_auth:0,black_list:"",white_list:"",clientip_flow_limit:"",service_flow_limit:""}}},methods:{onSubmit:function(){var e=this;this.submitButtonDisabled=!0;var t=Object.assign({},this.form);console.log(t),t.white_list=t.white_list.replace(/\n/g,","),t.black_list=t.black_list.replace(/\n/g,","),t.weight_list=t.weight_list.replace(/\n/g,","),t.ip_list=t.ip_list.replace(/\n/g,","),t.clientip_flow_limit=Number(t.clientip_flow_limit),t.service_flow_limit=Number(t.service_flow_limit),t.port=Number(t.port),this.isEdit?Object(l["j"])(t).then((function(t){e.submitButtonDisabled=!1,e.$notify({title:"Success",message:"Edit successfully!",type:"success",duration:2e3})})).catch((function(){e.submitButtonDisabled=!1})):Object(l["c"])(t).then((function(t){e.submitButtonDisabled=!1,e.$notify({title:"Success",message:"Add successfully!",type:"success",duration:2e3})})).catch((function(){e.submitButtonDisabled=!1}))},fetchData:function(e){var t=this,i={id:e};Object(l["e"])(i).then((function(e){t.form.id=e.data.info.id,t.form.load_type=e.data.info.load_type,t.form.service_name=e.data.info.service_name,t.form.port=e.data.tcp_rule.port.toString(),t.form.round_type=e.data.load_balance.round_type,t.form.ip_list=e.data.load_balance.ip_list.replace(/, /g,"\n"),t.form.weight_list=e.data.load_balance.weight_list.replace(/, /g,"\n"),t.form.open_auth=e.data.access_control.open_auth,t.form.black_list=e.data.access_control.black_list.replace(/, /g,"\n"),t.form.white_list=e.data.access_control.white_list.replace(/, /g,"\n"),t.form.clientip_flow_limit=e.data.access_control.clientip_flow_limit,t.form.service_flow_limit=e.data.access_control.service_flow_limit}))}},created:function(){_;var e=this.$route.params&&this.$route.params.id;e>0&&(this.isEdit=!0,this.fetchData(e))}},o=s,c=(i("37bc"),i("2877")),n=Object(c["a"])(o,r,a,!1,null,"0f4d538c",null);t["default"]=n.exports},a4d4:function(e,t,i){}}]);