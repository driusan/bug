var BugPage = React.createClass({
	render: function() {
		return (
		<div>
			<div className="col-md-8">
				<h2>{this.props.Title}</h2>
				<div>{this.props.Description}</div>
				<a href="#" onClick={this.props.onBack}>Return to list</a>
			</div>
			<div className="col-md-4">
				<BugList Title="Other Issues" 
					Bugs={this.props.AllBugs}
					onBugClicked={this.props.onOtherBugClicked}
				/>
			</div>
		</div>);
	}
});
