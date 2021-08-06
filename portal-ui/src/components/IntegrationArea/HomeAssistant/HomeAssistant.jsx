import { makeStyles } from '@material-ui/core/styles';
import React, { useState } from "react";
import classNames from 'classnames';
import ExpansionPanel from '@material-ui/core/ExpansionPanel';
import ExpansionPanelDetails from '@material-ui/core/ExpansionPanelDetails';
import ExpansionPanelSummary from '@material-ui/core/ExpansionPanelSummary';
import ExpansionPanelActions from '@material-ui/core/ExpansionPanelActions';
import InputLabel from '@material-ui/core/InputLabel';
import MenuItem from '@material-ui/core/MenuItem';
import FormControl from '@material-ui/core/FormControl';
import Select from '@material-ui/core/Select';
import Typography from '@material-ui/core/Typography';
import Chip from '@material-ui/core/Chip';
import Button from '@material-ui/core/Button';
import Divider from '@material-ui/core/Divider';
import { useMutation } from '@apollo/client';
import { selectDevicesForIntegration } from '../mutations/selectDevicesForIntegration';


const useStyles = makeStyles((theme) => ({
    root: {
        flexGrow: 1,
        padding: "1rem",
        margin: "3%",
        textAlign: 'center',
        backgroundColor: "rgb(199, 216, 223)",
        width: '90%',
    },
    heading: {
        fontSize: theme.typography.pxToRem(23),
    },
    secondaryHeading: {
        fontSize: theme.typography.pxToRem(20),
        color: theme.palette.text.secondary,
    },
    columnHeading: {
        fontSize: theme.typography.pxToRem(15),
        color: theme.palette.text.secondary,
    },
    deviceHeading: {
        fontSize: theme.typography.pxToRem(12),
        color: theme.palette.text.secondary,
        display: 'block',
    },
    icon: {
        verticalAlign: 'bottom',
        height: 20,
        width: 20,
    },
    details: {
        alignItems: 'center',
    },
    column: {
        flexBasis: '33.33%',
    },
    helper: {
        borderLeft: `2px solid ${theme.palette.divider}`,
        padding: `${theme.spacing.unit}px ${theme.spacing.unit * 2}px`,
    },
    link: {
        color: theme.palette.primary.main,
        textDecoration: 'none',
        '&:hover': {
            textDecoration: 'underline',
        },
    },
    deviceDetails: {
        display: 'block',
        marginTop: theme.spacing(2),
    },
    button: {
        display: 'block',
        marginTop: theme.spacing(2),
        backgroundColor: "#1682a3",
    },
    formControl: {
        margin: theme.spacing(1),
        minWidth: 400,
    },
    em: {
        fontWeight: "bold",
        fontStyle: "normal"
    },
}));



export function HomeAssistantIntegration({ integration }) {
    const classes = useStyles();

    const [selectedDevice, changeSelectedDevice] = useState("");
    const [selectDevicesForIntegrationMutation] = useMutation(selectDevicesForIntegration);
    const handleChange = (event) => {
        changeSelectedDevice(event.target.value);
    };

    const addToSelected = () => {
        integration.selectedDevices = [selectedDevice]
        devices = [selectedDevice.entityId];
        console.log({devices})
        // TODO: Add mutation here to add the entityID to the selected devices for home assistant
        selectDevicesForIntegrationMutation({
            variables: {
                integration: "hass",
                devices: devices,
            }
        })
        changeSelectedDevice("");
    }

    var devices = [...integration.value.devices].sort((a, b) => (a.entityId > b.entityId) ? 1 : ((b.entityId > a.entityId) ? -1 : 0));

    // Remove automations from the provided list
    devices = devices.filter(function (elem) {
        return !elem.entityId.includes("automation");
    });

    return (
        <ExpansionPanel >
            <ExpansionPanelSummary >
                <div className={classes.column}>
                    <Typography className={classes.heading}>{integration.name}</Typography>
                </div>
                <div className={classes.column}>
                    <Typography className={classes.secondaryHeading}>Select Devices</Typography>
                </div>
            </ExpansionPanelSummary>
            <ExpansionPanelDetails className={classes.details}>
                <div className={classes.column}>
                    <FormControl className={classes.formControl}>
                        <InputLabel id="demo-controlled-open-select-label">HomeAssistant Devices</InputLabel>
                        <Select
                            labelId="demo-controlled-open-select-label"
                            id="demo-controlled-open-select"
                            value={selectedDevice}
                            onChange={handleChange}
                        >
                            <MenuItem value="">
                                <em>None</em>
                            </MenuItem>
                            {
                                devices.map(device => <MenuItem value={device}>{device.entityId}</MenuItem>)
                            }
                        </Select>
                    </FormControl>
                    {selectedDevice !== "" ?
                        <div className={classes.deviceDetails}>
                            <Typography className={classes.deviceHeading} align="left"><em className={classes.em}>ENTITY ID:</em>     {selectedDevice.entityId}</Typography>
                            <Typography className={classes.deviceHeading} align="left"><em className={classes.em}>TYPE:</em>          {selectedDevice.type}</Typography>
                            <Typography className={classes.deviceHeading} align="left"><em className={classes.em}>AREA NAME:</em>     {selectedDevice.areaName}</Typography>
                            <Typography className={classes.deviceHeading} align="left"><em className={classes.em}>NAME:</em>          {selectedDevice.name}</Typography>
                            <Typography className={classes.deviceHeading} align="left"><em className={classes.em}>CURRENT STATE:</em> {selectedDevice.state}</Typography>
                            <Button className={classes.button} onClick={addToSelected}>Add</Button>
                        </div> : <div></div>}
                </div>
                <div className={classNames(classes.column, classes.helper)}>
                    <Typography className={classes.columnHeading}>Selected Devices</Typography>
                    {integration.value.selectedDevices.map(device =>
                        <Chip label={device.entityId} className={classes.chip} onDelete={() => { }} />
                    )}
                </div>
                <div className={classNames(classes.column, classes.helper)}>
                    <Typography variant="caption">
                        Select your destination of choice
                        <br />
                        <a href="#sub-labels-and-columns" className={classes.link}>
                            Learn more
                        </a>
                    </Typography>
                </div>
            </ExpansionPanelDetails>
            <Divider />
            <ExpansionPanelActions>
                <Button size="small">Cancel</Button>
                <Button size="small" color="primary">
                    Save
                </Button>
            </ExpansionPanelActions>
        </ExpansionPanel>
    )
}
