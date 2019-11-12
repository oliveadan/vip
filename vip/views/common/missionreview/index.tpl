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
                    <li class="layui-this">任务审核列表</li>
                </ul>
                <div class="layui-tab-content">
                    <form class="layui-form layui-form-pane form-container" style="max-width: 1500px;"
                          action='{{urlfor "IndexMissionReview.get"}}' method="get">
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
                            <label class="layui-form-label">任务列表</label>
                            <div class="layui-input-inline">
                                <select id="GameId" name="missionId" lay-verify="required" lay-filter="gameslt" lay-search>
                                    <option value="0">全部</option>
                                    {{range $i, $m := .missionList}}
                                        <option value="{{$m.Id}}" {{if eq $m.Id $.condArr.missionId}}selected="selected"{{end}} >{{$m.Describe}}</option>
                                    {{end}}
                                </select>
                            </div>
                        </div>
                        <div class="layui-inline">
                            <div class="layui-input-inline">
                                <select name="status">
                                    <option value="3" {{if eq 3 .condArr.status}}selected="selected"{{end}}>全部审核状态
                                    </option>
                                    <option value="1" {{if eq 1 .condArr.status}}selected="selected"{{end}}>审核通过
                                    </option>
                                    <option value="0" {{if eq 0 .condArr.status}}selected="selected"{{end}}>审核中
                                    </option>
                                    <option value="2" {{if eq 2 .condArr.status}}selected="selected"{{end}}>审核拒绝
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
                            <button class="layui-btn">搜索</button>
                            <a href='{{urlfor "IndexMissionReview.ReviewBatch"}}' class="layui-btn layui-btn-danger layui-btn-small ajax-batch">批量标记派送</a>
                            <!--
                                <a href='{{urlfor "RewardLogIndexController.Delbatch"}}'
                               class="layui-btn layui-btn-danger ajax-batch">删除所有</a>
                             -->
                        </div>
                    </form>
                    <hr>
                    <div class="layui-tab-item layui-show">
                        <table class="layui-table">
                            <thead>
                            <tr>
                                <th style="width:30px;">ID</th>
                                <th style="width:80px;">会员账号</th>
                                <th>活动描述</th>
                                <th>活动详情</th>
                                <th>最小VIP等级</th>
                                <th>最大VIP等级</th>
                                <th>积分</th>
                                <th>备注</th>
                                <th style="width:80px;">申请时间</th>
                                <th style="width:60px;">审核状态</th>
                                <th style="width:80px;">审核时间</th>
                                <th>操作</th>
                            </tr>
                            </thead>
                            <tbody>
                            {{range $index, $vo := .dataList}}
                                <tr>
                                    <td>{{$vo.Id}}</td>
                                    <td>{{$vo.Account}}</td>
                                    <td>{{getMissionDescribe $vo.MissionId}}</td>
                                    <td>{{getMissionDetail $vo.MissionDetailId}}</td>
                                    <td>{{$vo.MinLevel}}</td>
                                    <td>{{$vo.MaxLevel}}</td>
                                    <td>{{$vo.Integral}}</td>
                                    <td>{{$vo.Remark}}</td>
                                    <td>{{date $vo.CreateDate "Y-m-d H:i:s"}}</td>
                                    <td>{{if eq $vo.Status 0}}
                                            <span class="layui-badge layui-bg-orange">审核中</span>
                                        {{else if eq $vo.Status 1}}
                                            <span class="layui-badge layui-bg-green">审核通过</span>
                                        {{else if eq $vo.Status 2}}
                                            <span class="layui-badge layui-bg-red">审核拒绝</span>
                                        {{end}}
                                    </td>
                                    <td>{{date $vo.DeliveredTime "Y-m-d H:i:s"}}</td>
                                    <td>
                                        {{if eq $vo.Status 0}}
                                            <a href='{{urlfor "IndexMissionReview.Review" "id" $vo.Id "status" 1}}'
                                               class="layui-btn layui-btn-sm ajax-click">审核通过</a>
                                            <a href='{{urlfor "IndexMissionReview.Review" "id" $vo.Id "status" 2}}'
                                               class="layui-btn layui-btn-danger layui-btn-sm ajax-click">审核拒绝</a>
                                        {{else if eq $vo.Status 1}}
                                            <a href='{{urlfor "IndexMissionReview.Review" "id" $vo.Id "status" 2}}'
                                               class="layui-btn layui-btn-danger layui-btn-sm ajax-click">审核拒绝</a>
                                        {{else if eq $vo.Status 2}}
                                            <a href='{{urlfor "IndexMissionReview.Review" "id" $vo.Id "status" 1}}'
                                               class="layui-btn layui-btn-sm ajax-click">审核通过</a>
                                        {{end}}
                                        <button class="layui-btn layui-btn-sm" onclick="ChangeRemark({{$vo.Id}})">备注</button>
                                        <a href='{{urlfor "IndexMissionReview.Delone" "id" $vo.Id}}'
                                           class="layui-btn layui-btn-danger layui-btn-sm ajax-delete">删除</a>
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
<script>
    function ChangeRemark(id) {
        layer.open({
            type: 1,
            title: '备注',
            btn: ['提交', '取消'],
            content: '<div style="padding:30px;"><input type="text" style="outline-style:none; border: 1px solid #ccc;border-radius: 3px; padding: 13px 14px;width:270px;font-weight: 700; font-size: 14px; "  id="content" placeholder="请输入备注内容" ></div>',
            yes: function () {
                $.ajax({
                    url: {{urlfor "VipCenterController.Remark"}},
                    type: "post",
                    data: { "id":id,  "content": $("#content").val()},
                    success: function (info) {
                        if (info.code === 1) {
                            layer.closeAll();
                            setTimeout(function () {
                                location.href =location.href;
                                layer.msg(info.msg)
                            }, 1000);
                        } else {
                            $("#content").val("");
                        }
                        layer.msg(info.msg);
                    },
                });
            }
        });
    }
</script>
</body>
</html>