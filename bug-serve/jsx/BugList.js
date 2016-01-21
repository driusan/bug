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
			return (<li key={"BugListItem" + val}>
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
