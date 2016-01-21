var BugApp = React.createClass({
	componentDidMount: function() {
		var that = this;
		AjaxGet("/issues/", function(response) {
			that.setState({
				"Bugs" : JSON.parse(response)
			});
		});
        AjaxGet("/settings", function(response) {
			that.setState({
				"Settings" : JSON.parse(response)
			});
        })
	},
	getInitialState : function() {
		return {
            "Settings" : {},
			"Title" : "Open Issues",
			"Bugs": [],
			"SelectedBugJSON" : null
		}
	},
	selectBugHandler: function(e) {
		e.preventDefault();
		var bug = e.currentTarget.textContent;
		var that = this;
		AjaxGet("/issues/" + bug + "/", function(response) {
			that.setState({SelectedBug : JSON.parse(response)});
		});
	},
	resetSelected: function() {
		this.setState({ "SelectedBug" : null});
	},
	render: function() {
		var content;
		if(this.state.SelectedBug != null) {
			content = <BugPage 
                Title={this.state.SelectedBug.Title}
                Description={this.state.SelectedBug.Description} 
                Milestone={this.state.SelectedBug.Milestone} 
                Status={this.state.SelectedBug.Status}
                Priority={this.state.SelectedBug.Priority}
                Identifier={this.state.SelectedBug.Identifier}
                Tags={this.state.SelectedBug.Tags}
                onBack={this.resetSelected}
                AllBugs={this.state.Bugs}
                onOtherBugClicked={this.selectBugHandler}
                />
		} else {
			content = <BugList 
                Title={this.state.Title} 
                Bugs={this.state.Bugs} 
                onBugClicked={this.selectBugHandler}
                />
		}
		return (<div>
			<h1>Issues for: {this.state.Settings.Title}</h1>
			<div>
				{content}
			</div>
		</div>);
	}
});

var AjaxGet = function(url, callback) {
	var xmlhttp = new XMLHttpRequest();
	xmlhttp.onreadystatechange = function() {
		if (this.readyState === 4 && this.status == 200) {
			callback(this.responseText)
		}
	}
	xmlhttp.open("GET", url, true);
	xmlhttp.send()
}
