/**
 * 后台JS主入口
 */
layui.use(['layer','element','form','upload'], function(){
    var $ = layui.jquery,
        layer = layui.layer,
        element = layui.element,
        form = layui.form,
	    upload = layui.upload;

    var ajax = $.ajax;
    $.extend({
        ajax: function(url, options) {
            if (typeof url === 'object') {
                options = url;
                url = undefined;
            }
            options = options || {};
            url = options.url;
            var xsrftoken = $('meta[name=_xsrf]').attr('content');
            var headers = options.headers || {};
            var domain = document.domain.replace(/\./ig, '\\.');
            if (!/^(http:|https:).*/.test(url) || eval('/^(http:|https:)\\/\\/(.+\\.)*' + domain + '.*/').test(url)) {
                headers = $.extend(headers, {'X-Xsrftoken':xsrftoken});
            }
            options.headers = headers;
            return ajax(url, options);
        }
    });

    /**
     * 通用表单提交(AJAX方式)
     */
    form.on('submit(*)', function (data) {
        $.ajax({
            url: data.form.action,
            type: data.form.method,
            data: $(data.form).serialize(),
            success: function (info) {
                if (info.code === 1) {
                    setTimeout(function () {
                        location.href = info.url || location.href;
                    }, 1000);
                    layer.msg(info.msg, {icon: 1});
                } else {
                    layer.msg(info.msg, {icon: 2});
                }
            },
            error: function(info) {
                layer.msg(info.responseText || '请求异常', {icon: 2});
            }
        });

        return false;
    });

    /**
     * 通用批量处理（审核、取消审核、删除）
     */
    $('.ajax-action').on('click', function () {
        var _action = $(this).data('action');
        layer.open({
            shade: false,
            content: '确定执行此操作？',
            btn: ['确定', '取消'],
            yes: function (index) {
                $.ajax({
                    url: _action,
                    type: 'POST',
                    data: $('.ajax-form').serialize(),
                    success: function (info) {
                        if (info.code === 1) {
                            setTimeout(function () {
                                location.href = info.url || location.href;
                            }, 1000);
                            layer.msg(info.msg, {icon: 1});
                        } else {
                            layer.msg(info.msg, {icon: 2});
                        }
                    },
                    error: function(info) {
                        layer.msg(info.responseText || '请求异常', {icon: 2});
                    }
                });
                layer.close(index);
            }
        });

        return false;
    });

    /**
     * 通用全选
     */
    $('.check-all').on('click', function () {
        $(this).parents('table').find('input[type="checkbox"]').prop('checked', $(this).prop('checked'));
    });

    /**
     * 通用删除
     */
    $('.ajax-delete').on('click', function () {
        var _href = $(this).attr('href');
        layer.open({
            shade: false,
            content: '确定删除？',
            btn: ['确定', '取消'],
            yes: function (index) {
                $.ajax({
                    url: _href,
                    type: "POST",
                    success: function (info) {
                        if (info.code === 1) {
                            setTimeout(function () {
                                location.href = info.url || location.href;
                            }, 1000);
                            layer.msg(info.msg, {icon: 1});
                        } else {
                            layer.msg(info.msg, {icon: 2});
                        }
                    },
                    error: function(info) {
                        layer.msg(info.responseText || '请求异常', {icon: 2});
                    }
                });
                layer.close(index);
            }
        });

        return false;
    });

	/**
     * 通用批量提交操作
     */
    $('.ajax-batch').on('click', function () {
        var _href = $(this).attr('href');
        layer.open({
            shade: false,
            content: '该操作将根据条件批量处理数据，请谨慎操作！<br>确定执行？',
            btn: ['确定', '取消'],
            yes: function (index) {
                $.ajax({
                    url: _href,
                    type: "POST",
                    data: $('.layui-form').serialize(),
                    success: function (info) {
                        if (info.code === 1) {
                            setTimeout(function () {
                                location.href = info.url || location.href;
                            }, 1000);
                            layer.msg(info.msg, {icon: 1});
                        } else {
                            layer.msg(info.msg, {icon: 2});
                        }
                    },
                    error: function(info) {
                        layer.msg(info.responseText || '请求异常', {icon: 2});
                    }
                });
                layer.close(index);
            }
        });

        return false;
    });

    /**
     * 通用按钮点击事件
     */
    $('.ajax-click').on('click', function () {
        var _href = $(this).attr('href');
        layer.open({
            shade: false,
            content: '确定执行此操作？',
            btn: ['确定', '取消'],
            yes: function (index) {
                $.ajax({
                    url: _href,
                    type: "POST",
                    success: function (info) {
                        if (info.code === 1) {
                            setTimeout(function () {
                                location.href = info.url || location.href;
                            }, 1000);
                            layer.msg(info.msg, {icon: 1});
                        } else {
                            layer.msg(info.msg, {icon: 2});
                        }
                    },
                    error: function(info) {
                        layer.msg(info.responseText || '请求异常', {icon: 2});
                    }
                });
                layer.close(index);
            }
        });

        return false;
    });
	/*
	* 使用此方法，必须在元素上配置lay-data url
	*/
	upload.render({
	    elem: '#import',
		accept: 'file',
		before: function(obj){
	    	layer.load(); //上传loading
	  	},
	    done: function(res){
			layer.closeAll('loading');
			layer.open({
			  	title: '导入提示',
			  	content: res.msg
			});
	    },
	    error: function(){
			layer.closeAll('loading');
			layer.msg("导入失败，请刷新后重试");
	    }
  	});

    /**
     * 清除缓存
     */
    $('#clear-cache').on('click', function () {
        var _url = $(this).data('url');
        if (_url !== 'undefined') {
            $.ajax({
                url: _url,
                type: 'POST',
                success: function (data) {
                    if (data.code === 1) {
                        setTimeout(function () {
                            location.href = location.pathname;
                        }, 1000);
                        layer.msg(data.msg, {icon: 1});
                    } else {
                        layer.msg(data.msg, {icon: 2});
                    }
                }
            });
        }

        return false;
    });

});
