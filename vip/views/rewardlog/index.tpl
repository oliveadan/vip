<!DOCTYPE html>
<html lang="zh-CN">
<head>
    {{template "sysmanage/aalayout/meta.tpl" .}}
</head>
<body>
<div class="layui-fluid">
    <div class="layui-row layui-col-space10">
        <div class="layui-col-xs12 layui-col-sm12 layui-col-md12">
            <!--tab标签-->
            <div class="layui-tab layui-tab-brief">
                <ul class="layui-tab-title">
                    <li class="layui-this">中奖记录列表</li>
                </ul>
                <div class="layui-tab-content">
                    <form class="layui-form layui-form-pane form-container" style="max-width: 1500px;"
                          action='{{urlfor "RewardLogIndexController.get"}}' method="get">
                        <div class="layui-inline">
                            <div class="layui-input-inline">
                                <input type="text" name="account" value="{{.condArr.account}}" placeholder="会员账号"
                                       class="layui-input">
                            </div>
                        </div>
                        <div class="layui-inline">
                            <div class="layui-input-inline">
                                <input type="text" name="timeStart" value="{{.condArr.timeStart}}" placeholder="起始时间"
                                       class="layui-input">
                            </div>
                            <div class="layui-input-inline">
                                <input type="text" name="timeEnd" value="{{.condArr.timeEnd}}" placeholder="截止时间"
                                       class="layui-input">
                            </div>
                        </div>
                        <div class="layui-inline">
                            <div class="layui-input-inline">
                                <select name="delivered">
                                    <option value="2" {{if eq 2 .condArr.delivered}}selected="selected"{{end}}>全部派送状态
                                    </option>
                                    <option value="0" {{if eq 0 .condArr.delivered}}selected="selected"{{end}}>未派送
                                    </option>
                                    <option value="1" {{if eq 1 .condArr.delivered}}selected="selected"{{end}}>已派送
                                    </option>
                                </select>
                            </div>
                        </div>
                        <div class="layui-inline">
                            <div class="layui-input-inline">
                                <input type="checkbox" name="isExport" value="1" lay-skin="switch" lay-text="导出|导出">
                            </div>
                        </div>
                        <div class="layui-inline">
                            <button class="layui-btn">搜索/导出</button>
                            <a href='{{urlfor "RewardLogIndexController.Deliveredbatch"}}'
                               class="layui-btn layui-btn-small ajax-batch">批量标记派送</a>
                            <a href='{{urlfor "RewardLogIndexController.Delbatch"}}'
                               class="layui-btn layui-btn-danger ajax-batch">删除所有</a>
                        </div>
                    </form>
                    <hr>
                    <div class="layui-tab-item layui-show">
                        <table class="layui-table">
                            <thead>
                            <tr>
                                <th style="width:30px;">ID</th>
                                <th style="width:80px;">会员账号</th>
                                <th>奖品名称</th>
                                <th>奖品内容</th>
                                <th style="width:80px;">中奖时间</th>
                                <th style="width:60px;">是否派送</th>
                                <th style="width:80px;">派送时间</th>
                                <th>操作</th>
                            </tr>
                            </thead>
                            <tbody>
                            {{range $index, $vo := .dataList}}
                                <tr>
                                    <td>{{$vo.Id}}</td>
                                    <td>{{$vo.Account}}</td>
                                    <td>{{$vo.GiftName}}</td>
                                    <td>{{$vo.GiftContent}}</td>
                                    <td>{{date $vo.CreateDate "Y-m-d H:i:s"}}</td>
                                    <td>{{if eq $vo.Delivered 1}}<span
                                                class="layui-badge layui-bg-green">已派送</span>{{else}}<span
                                                class="layui-badge layui-bg-red">未派送</span>{{end}}</td>
                                    <td>{{date $vo.DeliveredTime "Y-m-d H:i:s"}}</td>
                                    <td>
                                        {{if eq $vo.Delivered 0}}
                                            <a href='{{urlfor "RewardLogIndexController.Delivered" "id" $vo.Id "delivered" 1}}'
                                               class="layui-btn layui-btn-mini ajax-click">标记派送</a>
                                        {{else}}
                                            <a href='{{urlfor "RewardLogIndexController.Delivered" "id" $vo.Id "delivered" 0}}'
                                               class="layui-btn layui-btn-primary layui-btn-mini ajax-click">取消派送</a>
                                        {{end}}

                                        <a href='{{urlfor "RewardLogIndexController.Delone" "id" $vo.Id}}'
                                           class="layui-btn layui-btn-danger layui-btn-mini ajax-delete">删除</a>
                                    </td>
                                </tr>
                            {{else}}
                                {{template "sysmanage/aalayout/table-no-data.tpl"}}
                            {{end}}
                            </tbody>
                        </table>
                        {{template "sysmanage/aalayout/paginator.tpl" .}}
                    </div>
                </div>
            </div>
        </div>
    </div>
</div>
{{template "sysmanage/aalayout/footjs.tpl" .}}
<script>
    layui.use('laydate', function () {
        var laydate = layui.laydate;
        laydate.render({
            elem: 'input[name="timeStart"]',
            type: 'datetime'
        });
        laydate.render({
            elem: 'input[name="timeEnd"]',
            type: 'datetime'
        });
    });
</script>
</body>
</html>