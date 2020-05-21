import React, { useState, useEffect } from 'react';
import { useHistory } from "react-router-dom";
// Externals
import PropTypes from 'prop-types';

// Material helpers
import { withStyles } from '@material-ui/core';

// Material components
import { Grid, CircularProgress } from '@material-ui/core';

// Shared layouts
import { Dashboard as DashboardLayout } from 'layouts';
import { useStoreState, useStoreActions } from 'easy-peasy';

/*
// Custom components
import RoomDetails from './components/RoomDetails';
import RoomGateways from './components/RoomGateways';
*/
import Detail from "./components/PropDetails";


import RoomDetail from '../../components/RoomPicker';
import CardDetail from '../../components/CardPicker';
import PropLocation from '../../components/PropLocation';

// Component styles
const styles = theme => ({
	root: {
		padding: theme.spacing.unit * 4
	}
});

const Property = (props) => {
	const { classes, match } = props;
	const { params } = match;
	const { id } = params;

	const history = useHistory();

	const baseUrl = "/zariadenia";
	const newIDName = "nove";

	const [inited, setInited] = useState(false);
	const [tabIndex, setTabIndex] = useState(0);


	const prop = useStoreState(state => state.props.prop);
	const propPending = useStoreState(state => state.props.propLoading);
	const propError = useStoreState(state => state.props.propError);
	const loadProp = useStoreActions(actions => actions.props.loadProp);
	const clearProp = useStoreActions(actions => actions.props.clearProp);
	const saveProp = useStoreActions(actions => actions.props.saveProp);
	const addItem = useStoreActions(actions => actions.props.addProp);


	const performLinkDelete = useStoreActions(actions => actions.props.deleteCardLink);

	const addFromLog = useStoreActions(actions => actions.props.addCardLinkFromLog);
	const fromLogAdding = useStoreState(state => state.props.cardLinkAdding);
	
	const loadCards = useStoreActions(actions => actions.cards.loadCard);

	const saveRoom = (room) => {
		if (room) {
			const editedProp = {
				...prop,
				room_id: room.id,
				room: null,
			}
			console.log("save prop", editedProp);
			saveProp(editedProp);
		}
	}

	const saveEditedProp = (prop) => {
		if (prop) {
			if (prop.price) {
				const prc = prop.price.replace && prop.price.replace(",", ".");
				prop.price = parseFloat(prc);
			}
			const editedProp = {
				...prop,
				name: prop.name,
				code: prop.code,
				price: prop.price,
				acquired_at: prop.acquired_at,
			}
			saveProp(editedProp);
		}
	}


	const addNewItem = (prop) => {
		if (prop) {
			if (prop.price) {
				const prc = prop.price.replace && prop.price.replace(",", ".");
				prop.price = parseFloat(prc);
			}
	
			addItem(prop).then(p => {
				console.log("created", p)
				history.push(`${baseUrl}/${p.id}`)
			});
		}
	}
	useEffect(() => {
		if (!inited) {
			setInited(true);
			if (id !== newIDName) {
				loadProp(params.id);
			} else {
				clearProp();
			}
		}
	}, []);

	useEffect(() => {
		console.log("prop", prop);
		console.log("prop pen", propPending);
	}, [prop]);

	const windowName = prop && prop.name ? prop.name : `Zariadenie`

	const emptyView = id == newIDName;
	console.log(propError);
	return (
		<DashboardLayout title={windowName}>
			<div className={classes.root}>
				{propPending &&
					<div className={classes.progressWrapper}>
						<CircularProgress />
					</div>
				}
				{!propPending && propError &&
					<div className={classes.progressWrapper}>
						<h2>Počas načítania nastala chyba</h2>
						{
							propError && propError.code === 404 &&
							<h3>
								Miestnosť sa nenašla.
							</h3>
						}
					</div>
				}
				{!propPending && !propError &&
					<>
						<Grid
							container
							spacing={4}
						>
							<Grid
								item
								lg={8}
								md={6}
								xl={8}
								xs={12}
							>
								<Detail
									data={prop}
									save={saveEditedProp}
									add={addNewItem}
									newID={emptyView}
								/>
							</Grid>
							<Grid
								item
								lg={4}
								md={6}
								xl={4}
								xs={12}
							>
								{!emptyView && <RoomDetail
									id={prop.room_id}
									save={saveRoom}
								/>}
							</Grid>
						</Grid>
						<Grid
							container
							spacing={4}
						>
							<Grid
								item
								lg={8}
								md={6}
								xl={8}
								xs={12}
							>
								{!emptyView && <PropLocation
									id={prop.id}
								/>}
							</Grid>
							<Grid
								item
								lg={4}
								md={6}
								xl={4}
								xs={12}
							>
								{!emptyView && <CardDetail
									id={prop.id}
									label="Informácie o RFID tagu"
									emptyText="Zariadenie nemá pridelený žiadny tag"
									performLinkDelete={performLinkDelete}
									addFromLog={addFromLog}
									fromLogAdding={fromLogAdding}
									loadCards={loadCards}
									controllerContext="props"
								/>}
							</Grid>
						</Grid>
					</>
				}
			</div>
		</DashboardLayout>
	);
}

Property.propTypes = {
	classes: PropTypes.object.isRequired
};

export default withStyles(styles)(Property);
