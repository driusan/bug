var BugPage = React.createClass({
    componentDidMount: function() {
        if(this.refs.desc) {
            this.refs.desc.innerHTML = marked(this.props.Description)
        }
        this.setState({
            "Editing" : false
        });
    },
    componentDidUpdate: function() {
        if(this.refs.desc) {
            this.refs.desc.innerHTML = marked(this.props.Description)
        }
    },
    loadPreviousBug: function() {
        for(var i = 1; i < this.props.AllBugs.length; i += 1) {
            var bugCandidate = this.props.AllBugs[i];
            if (bugCandidate === this.props.CurrentBug) {
                this.props.LoadBug(this.props.AllBugs[i-1]);
                return
            }
        }
        return;
    },
    loadNextBug: function() {
        for(var i = 0; i < this.props.AllBugs.length - 1 ; i += 1) {
            var bugCandidate = this.props.AllBugs[i];
            if (bugCandidate === this.props.CurrentBug) {
                this.props.LoadBug(this.props.AllBugs[i+1]);
                return
            }
        }
        return;
    },
    getInitialState: function() {
        return {
            "Editing" : false
        }
    },
    editCurrentBug: function() {
        this.setState({
            "Editing" : true
        });
    },
    cancelEditting: function() {
        this.setState({
            "Editing" : false
        });
    },
    saveCurrentBug: function(e) {
        var editClosure = function(ctx) {
            return function(e) {
                ctx.setState({"Editing" : false});
                ctx.props.LoadBug(ctx.props.CurrentBug);
            }
        }
        AjaxPut(
            "/issues/" + this.props.CurrentBug + "/Description",
            this.refs.editdesc.value,
            editClosure(this),
            this.props.CurrentETag
        );

    },
    deleteCurrentBug: function() {
        var that = this;
        AjaxDelete(
            "/issues/" + this.props.CurrentBug,
            function() {
                if (that.props.onDelete) {
                    that.props.onDelete();
                }
            },
            this.props.CurrentETag
        );
    },
    onOtherBugClicked: function(e) {
        this.setState({
            "Editing" : false
        });

        this.props.onOtherBugClicked(e)
    },
	render: function() {
        var fieldRow = function(name, value) {
            if (value) {
            return (
                <div className="row field">
                    <div className="col-md-1">&nbsp;</div>
                    <div className="col-md-2 label label-info">
                        {name}
                    </div>
                    <div className="col-md-4 badge">
                        {value}
                    </div>
                    <div className="col-md-5">&nbsp;</div>
                </div>);
            }
            return;
        }, priority = fieldRow("Priority", this.props.Priority),
           statusRow = fieldRow("Status", this.props.Status),
           milestone = fieldRow("Milestone", this.props.Milestone);
        var prevClass, nextClass;
        if (this.props.AllBugs.length > 1 && 
                this.props.AllBugs[0] != this.props.CurrentBug) {
            prevClass = "previous";
        } else {
            prevClass = "previous disabled";
        }
        if (this.props.AllBugs.length > 1 && 
                this.props.AllBugs[this.props.AllBugs.length-1] != this.props.CurrentBug) {
            nextClass = "next";
        } else {
            nextClass = "next disabled";
        }

        var descDiv;
        if (this.state.Editing === true) {
            descDiv = (<div>
                <div className="row">
                    <textarea className="col-md-12" rows="16" defaultValue={this.props.Description} ref="editdesc" />
                </div>
                <div className="row">
                <div className="col-md-10">
                    <button className="btn btn-default" onClick={this.saveCurrentBug}>Save</button>
                    <button className="btn btn-default" onClick={this.cancelEditting}>Cancel</button>
                </div>
                <div className="col-md-2">
                    <button onClick={this.deleteCurrentBug} className="btn btn-default glyphicon glyphicon-remove">Delete</button>
                    </div>
                </div>
            </div>);
        } else {
            descDiv = <div ref="desc" onClick={this.editCurrentBug}>{this.props.Description}</div>
        }
		return (
		<div>
			<div className="col-md-8 container">
                <div className="jumbotron bugsummary">
                    <h2>{this.props.Title}</h2>
                    {descDiv}
                </div>
                <div className="fields">
                    {priority}
                    {statusRow}
                    {milestone}
                </div>
			</div>
			<div className="col-md-4">
				<BugList Title="Other Issues" 
					Bugs={this.props.AllBugs}
					onBugClicked={this.onOtherBugClicked}
				/>
			</div>
            <div className="col-md-12">
                <ul className="pager">
                    <li className={prevClass}><a href="#" onClick={this.loadPreviousBug}><span aria-hidden="true">&larr;</span> Previous</a></li>
                    <li className="return"><a href="#" onClick={this.props.onBack}>Return to list</a></li>
                    <li className={nextClass}><a href="#" onClick={this.loadNextBug}>Next <span aria-hidden="true">&rarr;</span></a></li>
                </ul>
            </div>
		</div>);
	}
});
