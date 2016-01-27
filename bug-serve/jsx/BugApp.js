var BugApp = React.createClass({
    reloadIssues: function() {
        var that = this;
		AjaxGet("/issues/", function(response) {
			that.setState({
				"Bugs" : JSON.parse(response)
			});
		});
    },
    reloadSettings: function() {
        var that = this;
        AjaxGet("/settings", function(response) {
			that.setState({
				"Settings" : JSON.parse(response)
			});
        })
    },
	componentDidMount: function() {
        this.reloadIssues();
        this.reloadSettings();
	},
	getInitialState : function() {
		return {
            "Settings" : {},
			"Title" : "Open Issues",
			"Bugs": [],
			"SelectedBug" : null,
            "SelectedBugDir" : ""
		}
	},
    loadBug: function(bug) {
		var that = this;
		AjaxGet("/issues/" + bug + "/", function(response, handler) {
			that.setState({
                "SelectedBug" : JSON.parse(response),
                "SelectedBugETag" : handler.getResponseHeader("ETag"),
                "SelectedBugDir" : bug
            });
		});
    },
	selectBugHandler: function(e) {
		e.preventDefault();
        this.loadBug(e.currentTarget.textContent);
	},
	resetSelected: function() {
		this.setState({ "SelectedBug" : null,
            "SelectedBugETag" : null,
            "SelectedBugDir" : null
            });
	},
    resetBugs: function() {
        this.resetSelected();
        this.reloadIssues();
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
                onDelete={this.resetBugs}
                AllBugs={this.state.Bugs}
                CurrentBug={this.state.SelectedBugDir}
                CurrentETag={this.state.SelectedBugETag}
                onOtherBugClicked={this.selectBugHandler}
                LoadBug={this.loadBug}
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
			callback(this.responseText, this)
		}
	}
	xmlhttp.open("GET", url, true);
	xmlhttp.send()
}


var AjaxPut = function(url, data, callback, etag) {
	var xmlhttp = new XMLHttpRequest();
	xmlhttp.onreadystatechange = function() {
		if (this.readyState === 4 && (this.status >= 200 && this.status < 300)) {
			callback(this.responseText)
		}
	}
	xmlhttp.open("PUT", url, true);
    if (etag) {
        xmlhttp.setRequestHeader("If-Match", etag)
    }
	xmlhttp.send(data)
}

var AjaxPutJSON = function(url, data, callback) {
    AjaxPut(url, JSON.stringify(data), callback);
}

var AjaxDelete = function(url, callback, etag) {
	var xmlhttp = new XMLHttpRequest();
	xmlhttp.onreadystatechange = function() {
		if (this.readyState === 4 && (this.status >= 200 && this.status < 300)) {
			callback(this.responseText)
		}
	}
	xmlhttp.open("DELETE", url, true);
    if (etag) {
        xmlhttp.setRequestHeader("If-Match", etag)
    }
	xmlhttp.send()
}
