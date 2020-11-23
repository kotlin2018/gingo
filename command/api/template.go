package api

var feApi = `package fe

import service from '@/utils/request'

export const create{{.StructName}} = (data) => {
     return service({
         url: "/{{.Abbreviation}}/create{{.StructName}}",
         method: 'post',
         data
     })
 }

 export const delete{{.StructName}} = (data) => {
     return service({
         url: "/{{.Abbreviation}}/delete{{.StructName}}",
         method: 'delete',
         data
     })
 }

 export const delete{{.StructName}}ByIds = (data) => {
     return service({
         url: "/{{.Abbreviation}}/delete{{.StructName}}ByIds",
         method: 'delete',
         data
     })
 }

 export const update{{.StructName}} = (data) => {
     return service({
         url: "/{{.Abbreviation}}/update{{.StructName}}",
         method: 'put',
         data
     })
 }

 export const find{{.StructName}} = (params) => {
     return service({
         url: "/{{.Abbreviation}}/find{{.StructName}}",
         method: 'get',
         params
     })
 }

 export const get{{.StructName}}List = (params) => {
     return service({
         url: "/{{.Abbreviation}}/get{{.StructName}}List",
         method: 'get',
         params
     })
 }
`
var feTable = `package fe
<template>
  <div>
    <div class="search-term">
      <el-form :inline="true" :model="searchInfo" class="demo-form-inline">
           {{- range .Fields}}  {{- if .FieldSearchType}} {{- if eq .FieldType "bool" }}
            <el-form-item label="{{.FieldDesc}}" prop="{{.FieldJson}}">
            <el-select v-model="searchInfo.{{.FieldJson}}" clear placeholder="请选择">
                <el-option
                    key="true"
                    label="是"
                    value="true">
                </el-option>
                <el-option
                    key="false"
                    label="否"
                    value="false">
                </el-option>
            </el-select>
            </el-form-item>
                  {{- else }}
        <el-form-item label="{{.FieldDesc}}">
          <el-input placeholder="搜索条件" v-model="searchInfo.{{.FieldJson}}"></el-input>
        </el-form-item> {{ end }} {{ end }}  {{ end }}
        <el-form-item>
          <el-button @click="onSubmit" type="primary">查询</el-button>
        </el-form-item>
        <el-form-item>
          <el-button @click="openDialog" type="primary">新增{{.Description}}</el-button>
        </el-form-item>
        <el-form-item>
          <el-popover placement="top" v-model="deleteVisible" width="160">
            <p>确定要删除吗？</p>
              <div style="text-align: right; margin: 0">
                <el-button @click="deleteVisible = false" size="mini" type="text">取消</el-button>
                <el-button @click="onDelete" size="mini" type="primary">确定</el-button>
              </div>
            <el-button icon="el-icon-delete" size="mini" slot="reference" type="danger">批量删除</el-button>
          </el-popover>
        </el-form-item>
      </el-form>
    </div>
    <el-table
      :data="tableData"
      @selection-change="handleSelectionChange"
      border
      ref="multipleTable"
      stripe
      style="width: 100%"
      tooltip-effect="dark"
    >
    <el-table-column type="selection" width="55"></el-table-column>
    <el-table-column label="日期" width="180">
         <template slot-scope="scope">{{ "{{scope.row.CreatedAt|formatDate}}" }}</template>
    </el-table-column>
    {{range .Fields}}
    {{- if .DictType}}
      <el-table-column label="{{.FieldDesc}}" prop="{{.FieldJson}}" width="120">
        <template slot-scope="scope">
          {{"{{"}}filterDict(scope.row.{{.FieldJson}},"{{.DictType}}"){{"}}"}}
        </template>
      </el-table-column>
    {{- else if eq .FieldType "bool" }}
    <el-table-column label="{{.FieldDesc}}" prop="{{.FieldJson}}" width="120">
         <template slot-scope="scope">{{ "{{scope.row."}}{{.FieldJson}}{{"|formatBoolean}}" }}</template>
    </el-table-column> {{- else }}
    <el-table-column label="{{.FieldDesc}}" prop="{{.FieldJson}}" width="120"></el-table-column> {{ end }}
    {{ end }}
      <el-table-column label="按钮组">
        <template slot-scope="scope">
          <el-button @click="update{{.StructName}}(scope.row)" size="small" type="primary">变更</el-button>
          <el-popover placement="top" width="160" v-model="scope.row.visible">
            <p>确定要删除吗？</p>
            <div style="text-align: right; margin: 0">
              <el-button size="mini" type="text" @click="scope.row.visible = false">取消</el-button>
              <el-button type="primary" size="mini" @click="delete{{.StructName}}(scope.row)">确定</el-button>
            </div>
            <el-button type="danger" icon="el-icon-delete" size="mini" slot="reference">删除</el-button>
          </el-popover>
        </template>
      </el-table-column>
    </el-table>

    <el-pagination
      :current-page="page"
      :page-size="pageSize"
      :page-sizes="[10, 30, 50, 100]"
      :style="{float:'right',padding:'20px'}"
      :total="total"
      @current-change="handleCurrentChange"
      @size-change="handleSizeChange"
      layout="total, sizes, prev, pager, next, jumper"
    ></el-pagination>

    <el-dialog :before-close="closeDialog" :visible.sync="dialogFormVisible" title="弹窗操作">
      此处请使用表单生成器生成form填充 表单默认绑定 formData 如手动修改过请自行修改key
      <div class="dialog-footer" slot="footer">
        <el-button @click="closeDialog">取 消</el-button>
        <el-button @click="enterDialog" type="primary">确 定</el-button>
      </div>
    </el-dialog>
  </div>
</template>

<script>
import {
    create{{.StructName}},
    delete{{.StructName}},
    delete{{.StructName}}ByIds,
    update{{.StructName}},
    find{{.StructName}},
    get{{.StructName}}List
} from "@/api/{{.PackageName}}";  //  此处请自行替换地址
import { formatTimeToStr } from "@/utils/data";
import infoList from "@/components/mixins/infoList";

export default {
  name: "{{.StructName}}",
  mixins: [infoList],
  data() {
    return {
      listApi: get{{.StructName}}List,
      dialogFormVisible: false,
      visible: false,
      type: "",
      deleteVisible: false,
      multipleSelection: [],
      {{- range .Fields}}
          {{- if .DictType }}
            {{.DictType}}Options:[],
          {{ end -}}
      {{end -}}
      formData: {
        {{range .Fields}}{{.FieldJson}}:null,{{ end }}
      }
    };
  },
  filters: {
    formatDate: function(time) {
      if (time != null && time != "") {
        var date = new Date(time);
        return formatTimeToStr(date, "yyyy-MM-dd hh:mm:ss");
      } else {
        return "";
      }
    },
    formatBoolean: function(bool) {
      if (bool != null) {
        return bool ? "是" :"否";
      } else {
        return "";
      }
    }
  },
  methods: {
      //条件搜索前端看此方法
      onSubmit() {
        this.page = 1
        this.pageSize = 10
        {{- range .Fields}} {{- if eq .FieldType "bool" }}      
        if (this.searchInfo.{{.FieldJson}}==""){
          this.searchInfo.{{.FieldJson}}=null
        } {{ end }} {{ end }}    
        this.getTableData()
      },
      handleSelectionChange(val) {
        this.multipleSelection = val
      },
      async onDelete() {
        const ids = []
        this.multipleSelection &&
          this.multipleSelection.map(item => {
            ids.push(item.ID)
          })
        const res = await delete{{.StructName}}ByIds({ ids })
        if (res.code == 0) {
          this.$message({
            type: 'success',
            message: '删除成功'
          })
          this.deleteVisible = false
          this.getTableData()
        }
      },
    async update{{.StructName}}(row) {
      const res = await find{{.StructName}}({ ID: row.ID });
      this.type = "update";
      if (res.code == 0) {
        this.formData = res.data.re{{.Abbreviation}};
        this.dialogFormVisible = true;
      }
    },
    closeDialog() {
      this.dialogFormVisible = false;
      this.formData = {
        {{range .Fields}}
          {{.FieldJson}}:null,{{ end }}
      };
    },
    async delete{{.StructName}}(row) {
      this.visible = false;
      const res = await delete{{.StructName}}({ ID: row.ID });
      if (res.code == 0) {
        this.$message({
          type: "success",
          message: "删除成功"
        });
        this.getTableData();
      }
    },
    async enterDialog() {
      let res;
      switch (this.type) {
        case "create":
          res = await create{{.StructName}}(this.formData);
          break;
        case "update":
          res = await update{{.StructName}}(this.formData);
          break;
        default:
          res = await create{{.StructName}}(this.formData);
          break;
      }
      if (res.code == 0) {
        this.$message({
          type:"success",
          message:"创建/更改成功"
        })
        this.closeDialog();
        this.getTableData();
      }
    },
    openDialog() {
      this.type = "create";
      this.dialogFormVisible = true;
    }
  },
  async created() {
    await this.getTableData();
  {{- range .Fields -}}
    {{- if .DictType -}}
      await this.getDict("{{.DictType}}")
    {{- end -}}
  {{- end -}}
}
};
</script>

<style>
</style>
`

var apiGo = `package v1

import (
	"fmt"
	"{{.Appname}}/global/response"
	"{{.Appname}}/model"
	"{{.Appname}}/model/request"
	resp "{{.Appname}}/model/response"
	"{{.Appname}}/service"
	"github.com/gin-gonic/gin"
)

func Create{{.StructName}}(c *gin.Context) {
	var {{.Abbreviation}} model.{{.StructName}}
	_ = c.ShouldBindJSON(&{{.Abbreviation}})
	err := service.Create{{.StructName}}({{.Abbreviation}})
	if err != nil {
		response.FailWithMessage(fmt.Sprintf("创建失败，%v", err), c)
	} else {
		response.OkWithMessage("创建成功", c)
	}
}

func Delete{{.StructName}}(c *gin.Context) {
	var {{.Abbreviation}} model.{{.StructName}}
	_ = c.ShouldBindJSON(&{{.Abbreviation}})
	err := service.Delete{{.StructName}}({{.Abbreviation}})
	if err != nil {
		response.FailWithMessage(fmt.Sprintf("删除失败，%v", err), c)
	} else {
		response.OkWithMessage("删除成功", c)
	}
}

func Delete{{.StructName}}ByIds(c *gin.Context) {
	var IDS request.IdsReq
    _ = c.ShouldBindJSON(&IDS)
	err := service.Delete{{.StructName}}ByIds(IDS)
	if err != nil {
		response.FailWithMessage(fmt.Sprintf("删除失败，%v", err), c)
	} else {
		response.OkWithMessage("删除成功", c)
	}
}

func Update{{.StructName}}(c *gin.Context) {
	var {{.Abbreviation}} model.{{.StructName}}
	_ = c.ShouldBindJSON(&{{.Abbreviation}})
	err := service.Update{{.StructName}}(&{{.Abbreviation}})
	if err != nil {
		response.FailWithMessage(fmt.Sprintf("更新失败，%v", err), c)
	} else {
		response.OkWithMessage("更新成功", c)
	}
}

func Find{{.StructName}}(c *gin.Context) {
	var {{.Abbreviation}} model.{{.StructName}}
	_ = c.ShouldBindQuery(&{{.Abbreviation}})
	err, re{{.Abbreviation}} := service.Get{{.StructName}}({{.Abbreviation}}.ID)
	if err != nil {
		response.FailWithMessage(fmt.Sprintf("查询失败，%v", err), c)
	} else {
		response.OkWithData(gin.H{"re{{.Abbreviation}}": re{{.Abbreviation}}}, c)
	}
}

func Get{{.StructName}}List(c *gin.Context) {
	var pageInfo request.{{.StructName}}Search
	_ = c.ShouldBindQuery(&pageInfo)
	err, list, total := service.Get{{.StructName}}InfoList(pageInfo)
	if err != nil {
		response.FailWithMessage(fmt.Sprintf("获取数据失败，%v", err), c)
	} else {
		response.OkWithData(resp.PageResult{
			List:     list,
			Total:    total,
			Page:     pageInfo.Page,
			PageSize: pageInfo.PageSize,
		}, c)
	}
}
`

var modelGo = `package model
// 自动生成模板{{.StructName}}

import (
	"github.com/jinzhu/gorm"
)

// 如果含有time.Time 请自行import time包
type {{.StructName}} struct {
      gorm.Model {{- range .Fields}}
            {{- if eq .FieldType "bool" }}
      {{.FieldName}}  *{{.FieldType}} // json:"{{.FieldJson}}" form:"{{.FieldJson}}" gorm:"column:{{.ColumnName}};comment:'{{.Comment}}'{{- if .DataType -}};type:{{.DataType}}{{- if .DataTypeLong -}}({{.DataTypeLong}}){{- end -}}{{- end -}}"
            {{- else }}
      {{.FieldName}}  {{.FieldType}} // json:"{{.FieldJson}}" form:"{{.FieldJson}}" gorm:"column:{{.ColumnName}};comment:'{{.Comment}}'{{- if .DataType -}};type:{{.DataType}}{{- if .DataTypeLong -}}({{.DataTypeLong}}){{- end -}}{{- end -}}"
            {{- end }}  {{- end }} 
}

{{ if .TableName }}
func ({{.StructName}}) TableName() string {
  return "{{.TableName}}"
}
{{ end }}
`

var requestGo =`package request

import "{{.Appname}}/model"

type {{.StructName}}Search struct{
    model.{{.StructName}}
    PageInfo
}
`

var routerGo = `package router

import (
	"{{.Appname}}/api/v1"
	"{{.Appname}}/middleware"
	"github.com/gin-gonic/gin"
)

func Init{{.StructName}}Router(Router *gin.RouterGroup) {
	{{.StructName}}Router := Router.Group("{{.Abbreviation}}").Use(middleware.JWTAuth()).Use(middleware.CasbinHandler()).Use(middleware.OperationRecord())
	{
		{{.StructName}}Router.POST("create{{.StructName}}", v1.Create{{.StructName}})   // 新建{{.StructName}}
		{{.StructName}}Router.DELETE("delete{{.StructName}}", v1.Delete{{.StructName}}) // 删除{{.StructName}}
		{{.StructName}}Router.DELETE("delete{{.StructName}}ByIds", v1.Delete{{.StructName}}ByIds) // 批量删除{{.StructName}}
		{{.StructName}}Router.PUT("update{{.StructName}}", v1.Update{{.StructName}})    // 更新{{.StructName}}
		{{.StructName}}Router.GET("find{{.StructName}}", v1.Find{{.StructName}})        // 根据ID获取{{.StructName}}
		{{.StructName}}Router.GET("get{{.StructName}}List", v1.Get{{.StructName}}List)  // 获取{{.StructName}}列表
	}
}
`

var serviceGo = `package service

import (
	"{{.Appname}}/global"
	"{{.Appname}}/model"
	"{{.Appname}}/model/request"
)

func Create{{.StructName}}({{.Abbreviation}} model.{{.StructName}}) (err error) {
	err = global.GVA_DB.Create(&{{.Abbreviation}}).Error
	return err
}

func Delete{{.StructName}}({{.Abbreviation}} model.{{.StructName}}) (err error) {
	err = global.GVA_DB.Delete({{.Abbreviation}}).Error
	return err
}

func Delete{{.StructName}}ByIds(ids request.IdsReq) (err error) {
	err = global.GVA_DB.Delete(&[]model.{{.StructName}}{},"id in (?)",ids.Ids).Error
	return err
}

func Update{{.StructName}}({{.Abbreviation}} *model.{{.StructName}}) (err error) {
	err = global.GVA_DB.Save({{.Abbreviation}}).Error
	return err
}

func Get{{.StructName}}(id uint) (err error, {{.Abbreviation}} model.{{.StructName}}) {
	err = global.GVA_DB.Where("id = ?", id).First(&{{.Abbreviation}}).Error
	return
}

func Get{{.StructName}}InfoList(info request.{{.StructName}}Search) (err error, list interface{}, total int) {
	limit := info.PageSize
	offset := info.PageSize * (info.Page - 1)
    // 创建db
	db := global.GVA_DB.Model(&model.{{.StructName}}{})
    var {{.Abbreviation}}s []model.{{.StructName}}
    // 如果有条件搜索 下方会自动创建搜索语句
        {{- range .Fields}}
            {{- if .FieldSearchType}}
                {{- if eq .FieldType "string" }}
    if info.{{.FieldName}} != "" {
        //db = db.Where("{{.ColumnName}} {{.FieldSearchType}} ?",{{if eq .FieldSearchType "LIKE"}}"%"+ {{ end }}info.{{.FieldName}}{{if eq .FieldSearchType "LIKE"}}+"%"{{ end }})
    }
                {{- else if eq .FieldType "bool" }}
    if info.{{.FieldName}} != nil {
       // db = db.Where("{{.ColumnName}} {{.FieldSearchType}} ?",{{if eq .FieldSearchType "LIKE"}}"%"+{{ end }}info.{{.FieldName}}{{if eq .FieldSearchType "LIKE"}}+"%"{{ end }})
    }
                {{- else if eq .FieldType "int" }}
    if info.{{.FieldName}} != 0 {
       // db = db.Where("{{.ColumnName}} {{.FieldSearchType}} ?",{{if eq .FieldSearchType "LIKE"}}"%"+{{ end }}info.{{.FieldName}}{{if eq .FieldSearchType "LIKE"}}+"%"{{ end }})
    }
                {{- else if eq .FieldType "float64" }}
    if info.{{.FieldName}} != 0 {
       // db = db.Where("{{.ColumnName}} {{.FieldSearchType}} ?",{{if eq .FieldSearchType "LIKE"}}"%"+{{ end }}info.{{.FieldName}}{{if eq .FieldSearchType "LIKE"}}+"%"{{ end }})
    }
                {{- else if eq .FieldType "time.Time" }}
    if !info.{{.FieldName}}.IsZero() {
        // db = db.Where("{{.ColumnName}} {{.FieldSearchType}} ?",{{if eq .FieldSearchType "LIKE"}}"%"+{{ end }}info.{{.FieldName}}{{if eq .FieldSearchType "LIKE"}}+"%"{{ end }})
    }
                {{- end }}
        {{- end }}
    {{- end }}
	err = db.Count(&total).Error
	err = db.Limit(limit).Offset(offset).Find(&{{.Abbreviation}}s).Error
	return err, {{.Abbreviation}}s, total
}
`