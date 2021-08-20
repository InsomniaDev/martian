import { makeStyles } from '@material-ui/core/styles';
import React, { useState } from "react";
import classNames from 'classnames';
import ExpansionPanel from '@material-ui/core/ExpansionPanel';
import ExpansionPanelDetails from '@material-ui/core/ExpansionPanelDetails';
import ExpansionPanelSummary from '@material-ui/core/ExpansionPanelSummary';
import InputLabel from '@material-ui/core/InputLabel';
import MenuItem from '@material-ui/core/MenuItem';
import FormControl from '@material-ui/core/FormControl';
import Select from '@material-ui/core/Select';
import Typography from '@material-ui/core/Typography';
import Chip from '@material-ui/core/Chip';
import Button from '@material-ui/core/Button';
import { useMutation } from '@apollo/client';
import { selectDevicesForIntegration } from '../mutations/selectDevicesForIntegration';
import HomeAssistantEditMenu from './HomeAssistantEditMenu';


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
    buttonDiv: {
        display: 'flex',
        flexDirection: 'row',
    },
    button: {
        display: 'block',
        margin: theme.spacing(1),
        backgroundColor: "#9cdbee",
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



export function HomeAssistantIntegration({ integration, refetchData, areaData }) {
    const classes = useStyles();

    const [selectedDevice, changeSelectedDevice] = useState("");
    const [selectDevicesForIntegrationMutation] = useMutation(selectDevicesForIntegration);
    const handleChange = (event) => {
        changeSelectedDevice(event.target.value);
    };

    // Add the selected variable to the interfaceDevices for Hass
    const addToSelectedInterface = () => {
        devices = [selectedDevice.entityId];
        selectDevicesForIntegrationMutation({
            variables: {
                integration: "hass",
                devices: devices,
                addDevices: true,
                automationDevice: false,
            }
        })
        refetchData();
    }

    // Remove the selected variable from the interfaceDevices for Hass
    const removeSelectedInterface = (device) => {
        selectDevicesForIntegrationMutation({
            variables: {
                integration: "hass",
                devices: [device.entityId],
                addDevices: false,
                automationDevice: false,
            }
        })
        refetchData();
    }

    // Add the selected variable to the interfaceDevices for Hass
    const addToSelectedAutomation = () => {
        devices = [selectedDevice.entityId];
        selectDevicesForIntegrationMutation({
            variables: {
                integration: "hass",
                devices: devices,
                addDevices: true,
                automationDevice: true,
            }
        })
        refetchData();
    }

    // Remove the selected variable from the interfaceDevices for Hass
    const removeSelectedAutomation = (device) => {
        selectDevicesForIntegrationMutation({
            variables: {
                integration: "hass",
                devices: [device.entityId],
                addDevices: false,
                automationDevice: true,
            }
        })
        refetchData();
    }

    const clearSelected = () => {
        changeSelectedDevice("");
    }

    const updateSelected = (newDevice) => {
        changeSelectedDevice(newDevice);
    }

    var devices = [...integration.value.devices].sort((a, b) => (a.entityId > b.entityId) ? 1 : ((b.entityId > a.entityId) ? -1 : 0));

    // Remove automations from the provided list
    devices = devices.filter(function (elem) {
        return !elem.entityId.includes("automation");
    });

    return (
        <ExpansionPanel key="hassExpansionPanel">
            <ExpansionPanelSummary key="hass" >
                <div key="hassTypographyNameDiv" className={classes.column}>
                    <Typography key="hassTypographyName" className={classes.heading}>{integration.name}</Typography>
                </div>
                <div key="hassTypographyHeadingDiv" className={classes.column}>
                    <Typography key="hassTypographyHeading" className={classes.secondaryHeading}>Edit Configuration</Typography>
                </div>
            </ExpansionPanelSummary>
            <ExpansionPanelDetails key="hassExpansionPanelDetails" className={classes.details}>
                <div key="hassExpansionPanelDetailsDiv" className={classes.column}>
                    <FormControl key="hassFormControl" className={classes.formControl}>
                        <InputLabel key="hassFormControlInputLabel" id="demo-controlled-open-select-label">HomeAssistant Devices</InputLabel>
                        <Select
                            key="hassFormControlSelect"
                            id="demo-controlled-open-select"
                            value={selectedDevice}
                            onChange={handleChange}
                        >
                            <MenuItem key="hassFormControlMenuItem" value="">
                                <em>None</em>
                            </MenuItem>
                            {
                                devices.map(device => <MenuItem value={device}>{device.entityId}</MenuItem>)
                            }
                        </Select>
                    </FormControl>
                    {selectedDevice !== "" ?
                        <div className={classes.deviceDetails}>
                            <Typography
                                key="selectedDeviceEntityId"
                                className={classes.deviceHeading}
                                align="left">
                                <em className={classes.em}>ENTITY ID:</em>     {selectedDevice.entityId}
                            </Typography>
                            <Typography
                                key="selectedDeviceType"
                                className={classes.deviceHeading}
                                align="left">
                                <em className={classes.em}>TYPE:</em>          {selectedDevice.type}
                            </Typography>
                            <Typography
                                key="selectedDeviceAreaName"
                                className={classes.deviceHeading}
                                align="left">
                                <em className={classes.em}>AREA NAME:</em>     {selectedDevice.areaName}
                            </Typography>
                            <Typography
                                key="selectedDeviceName"
                                className={classes.deviceHeading}
                                align="left">
                                <em className={classes.em}>NAME:</em>          {selectedDevice.name}
                            </Typography>
                            <Typography className={classes.deviceHeading} align="left"><em className={classes.em}>CURRENT STATE:</em> {selectedDevice.state}</Typography>
                            <div className={classes.buttonDiv}>
                                <Button className={classes.button} onClick={addToSelectedInterface}>Add to interface</Button>
                                <Button className={classes.button} onClick={addToSelectedAutomation}>Add to automated</Button>
                                <HomeAssistantEditMenu
                                    device={selectedDevice}
                                    buttonStyle={classes.button}
                                    buttonText="edit device"
                                    areaData={areaData}
                                    refetchData={refetchData}
                                    updateSelected={updateSelected} />
                                <Button className={classes.button} onClick={clearSelected}>clear</Button>
                            </div>
                        </div> : <div></div>}
                </div>
                <div className={classNames(classes.column, classes.helper)}>
                    <Typography className={classes.columnHeading}>Interface Devices</Typography>
                    {integration.value.interfaceDevices.map(device =>
                        <Chip label={device.entityId} className={classes.chip} onDelete={() => removeSelectedInterface(device)} />
                    )}
                </div>
                <div className={classNames(classes.column, classes.helper)}>
                    <Typography className={classes.columnHeading}>Automation Devices</Typography>
                    {integration.value.automatedDevices.map(device =>
                        <Chip label={device.entityId} className={classes.chip} onDelete={() => removeSelectedAutomation(device)} />
                    )}
                </div>
            </ExpansionPanelDetails>
        </ExpansionPanel>
    )
}
