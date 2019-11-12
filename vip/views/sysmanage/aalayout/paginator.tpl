<div>
	<ul class="pagination">
		<li class="disabled"><a>总记录数：{{.paginator.Nums}} 条</a></li>
	{{if .paginator.HasPrev}}
	    <li><a href="{{.paginator.PageLinkFirst}}">第一页</a></li>
	    <li><a href="{{.paginator.PageLinkPrev}}">&laquo;</a></li>
	{{else}}
	    <li class="disabled"><a>第一页</a></li>
	    <li class="disabled"><a>&laquo;</a></li>
	{{end}}
	{{range $index, $page := .paginator.Pages}}
	    <li{{if $.paginator.IsActive .}} class="active"{{end}}>
	        <a href="{{$.paginator.PageLink $page}}">{{$page}}</a>
	    </li>
	{{end}}
	{{if .paginator.HasNext}}
	    <li><a href="{{.paginator.PageLinkNext}}">&raquo;</a></li>
	    <li><a href="{{.paginator.PageLinkLast}}">最后一页</a></li>
	{{else}}
	    <li class="disabled"><a>&raquo;</a></li>
	    <li class="disabled"><a>最后一页</a></li>
	{{end}}
	</ul>
</div>