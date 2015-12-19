var BugList = React.createClass({
	getDefaultProps: function() {
		return {
			"Title" : "Bugs",
			"Bugs" : [],
			onBugClicked: function(e) { e.preventDefault(); return }
		}
	},
	render: function() {
		var that = this;
		var elements = 	this.props.Bugs.map(function (val) {
			return (<li>
				<a href="#" onClick={that.props.onBugClicked}>{val}</a>
				</li>);
		});
		return (<div>
			<h2>{this.props.Title}</h2>
			<ol>
				{elements}	
			</ol>
		</div>
		);
	}
});
/*
func (b BugListRenderer) GetBody() string {
	issues, _ := ioutil.ReadDir(bugs.GetRootDir() + "/issues")

	ret := "<h2>" + b.Title + "</h2><ol>"
	for _, issue := range issues {
		var dir bugs.Directory = bugs.Directory(issue.Name())
		ret += fmt.Sprintf("<li><a href=\"/issues/%s\">%s</a></li>\n", (dir), dir.ToTitle())
	}
	ret += "</ol>"

	return ret
}*/
