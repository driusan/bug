var BugPage = React.createClass({
    componentDidMount: function() {
        this.refs.desc.innerHTML = marked(this.props.Description)
    },
    componentDidUpdate: function() {
        this.refs.desc.innerHTML = marked(this.props.Description)
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
		return (
		<div>
			<div className="col-md-8 container">
                <div className="jumbotron bugsummary">
                    <h2>{this.props.Title}</h2>
                    <div ref="desc">{this.props.Description}</div>
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
					onBugClicked={this.props.onOtherBugClicked}
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
