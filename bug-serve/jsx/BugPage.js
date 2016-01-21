var BugPage = React.createClass({
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
                    <div classname="col-md-5">&nbsp;</div>
                </div>);
            }
            return;
        }, priority = fieldRow("Priority", this.props.Priority),
           statusRow = fieldRow("Status", this.props.Status),
           milestone = fieldRow("Milestone", this.props.Milestone);
		return (
		<div>
			<div className="col-md-8 container">
                <div className="jumbotron bugsummary">
                    <h2>{this.props.Title}</h2>
                    <div>{this.props.Description}</div>
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
                    <li className="previous"><a href="#"><span aria-hidden="true">&larr;</span> Previous</a></li>
                    <li className="return"><a href="#" onClick={this.props.onBack}>Return to list</a></li>
                    <li className="next"><a href="#">Next <span aria-hidden="true">&rarr;</span></a></li>
                </ul>
            </div>
		</div>);
	}
});
