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
import LutronEditMenu from './LutronEditMenu';


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
        padding: `${theme.spacing()}px ${theme.spacing(2)}px`,
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



export function LutronIntegration({ integration, refetchData, areaData }) {
    const classes = useStyles();

    const [selectedDevice, changeSelectedDevice] = useState("");
    const [selectDevicesForIntegrationMutation] = useMutation(selectDevicesForIntegration);
    const handleChange = (event) => {
        changeSelectedDevice(event.target.value);
    };

    // Add the selected variable to the interfaceDevices for lutron
    const addToSelectedInterface = () => {
        devices = [selectedDevice.id];
        selectDevicesForIntegrationMutation({
            variables: {
                integration: "lutron",
                devices: devices,
                addDevices: true,
                automationDevice: false,
            }
        })
        refetchData();
    }

    // Remove the selected variable from the interfaceDevices for lutron
    const removeSelectedInterface = (ipAddress) => {
        selectDevicesForIntegrationMutation({
            variables: {
                integration: "lutron",
                devices: [ipAddress],
                addDevices: false,
                automationDevice: false,
            }
        })
        refetchData();
    }

    // Add the selected variable to the interfaceDevices for lutron
    const addToSelectedAutomation = () => {
        devices = [selectedDevice.id];
        selectDevicesForIntegrationMutation({
            variables: {
                integration: "lutron",
                devices: devices,
                addDevices: true,
                automationDevice: true,
            }
        })
        refetchData();
    }

    // Remove the selected variable from the interfaceDevices for lutron
    const removeSelectedAutomation = (ipAddress) => {
        selectDevicesForIntegrationMutation({
            variables: {
                integration: "lutron",
                devices: [ipAddress],
                addDevices: false,
                automationDevice: true,
            }
        })
        refetchData();
    }

    const getNameForId = (id) => {
        const device = [...integration.value.devices].filter(device => device.id === id);
        return device[0].name;
    }

    const clearSelected = () => {
        changeSelectedDevice("");
    }

    const updateSelected = (newDevice) => {
        changeSelectedDevice(newDevice);
    }

    var devices = [...integration.value.devices].sort((a, b) => (a.lutronName > b.lutronName) ? 1 : ((b.lutronName > a.lutronName) ? -1 : 0));
    var interfaceDevices = [...integration.value.interfaceDevices].map(device => getNameForId(device)).sort((a, b) => (a > b) ? 1 : ((b > a) ? -1 : 0));
    var automatedDevices = [...integration.value.automatedDevices].map(device => getNameForId(device)).sort((a, b) => (a > b) ? 1 : ((b > a) ? -1 : 0));

    return (
        <ExpansionPanel key="lutronExpansionPanel">
            <ExpansionPanelSummary key="lutron" >
                <div key="lutronTypographyNameDiv" className={classes.column}>
                    <Typography key="lutronTypographyName" className={classes.heading}>{integration.name}</Typography>
                </div>
                <div key="lutronTypographyHeadingDiv" className={classes.column}>
                    <Typography key="lutronTypographyHeading" className={classes.secondaryHeading}>Edit Configuration</Typography>
                </div>
            </ExpansionPanelSummary>
            <ExpansionPanelDetails key="lutronExpansionPanelDetails" className={classes.details}>
                <div key="lutronExpansionPanelDetailsDiv" className={classes.column}>
                    <FormControl key="lutronFormControl" className={classes.formControl}>
                        <InputLabel key="lutronFormControlInputLabel" id="demo-controlled-open-select-label">Lutron Devices</InputLabel>
                        <Select
                            key="lutronFormControlSelect"
                            id="demo-controlled-open-select"
                            value={selectedDevice}
                            onChange={handleChange}
                        >
                            <MenuItem key="lutronFormControlMenuItem" value="">
                                <em>None</em>
                            </MenuItem>
                            {
                                devices.map(device => <MenuItem key={"select_" + device.areaName + device.name} value={device}><em>{device.areaName}</em>{". - " + device.name}</MenuItem>)
                            }
                        </Select>
                    </FormControl>
                    {selectedDevice !== "" ?
                        <div className={classes.deviceDetails}>
                            <Typography
                                key={"selectedDeviceIpAddress" + selectedDevice.id}
                                className={classes.deviceHeading}
                                align="left">
                                <em className={classes.em}>ID:</em>     {selectedDevice.id}
                            </Typography>
                            <Typography
                                key={"selectedDeviceType" + selectedDevice.id}
                                className={classes.deviceHeading}
                                align="left">
                                <em className={classes.em}>TYPE:</em>          {selectedDevice.type}
                            </Typography>
                            <Typography
                                key={"selectedDeviceAreaName" + selectedDevice.id}
                                className={classes.deviceHeading}
                                align="left">
                                <em className={classes.em}>AREA NAME:</em>     {selectedDevice.areaName}
                            </Typography>
                            <Typography
                                key={"selectedDeviceName" + selectedDevice.id}
                                className={classes.deviceHeading}
                                align="left">
                                <em className={classes.em}>NAME:</em>          {selectedDevice.name}
                            </Typography>
                            <div key={"selectedDeviceDiv" + selectedDevice.id} className={classes.buttonDiv}>
                                <Button key={"selectedDeviceInterfaceButton" + selectedDevice.id} className={classes.button} onClick={addToSelectedInterface}>Add to interface</Button>
                                <Button key={"selectedDeviceAutomatedButton" + selectedDevice.id} className={classes.button} onClick={addToSelectedAutomation}>Add to automated</Button>
                                <LutronEditMenu
                                    key={"selectedDeviceEditButton" + selectedDevice.id}
                                    device={selectedDevice}
                                    buttonStyle={classes.button}
                                    buttonText="edit device"
                                    areaData={areaData}
                                    refetchData={refetchData}
                                    updateSelected={updateSelected} />
                                <Button key={"selectedDeviceClearButton" + selectedDevice.id} className={classes.button} onClick={clearSelected}>clear</Button>
                            </div>
                        </div> : <div></div>}
                </div>
                <div key={"lutronInterfaceSelection" + selectedDevice.id} className={classNames(classes.column, classes.helper)}>
                    <Typography key={"lutronInterfaceSelectionTypography" + selectedDevice.id} className={classes.columnHeading}>Interface Devices</Typography>
                    {interfaceDevices.map(device =>
                        <Chip key={"interface_chip_" + device} label={device} className={classes.chip} onDelete={() => removeSelectedInterface(device)} />
                    )}
                </div>
                <div key={"lutronAutomationSelection" + selectedDevice.id} className={classNames(classes.column, classes.helper)}>
                    <Typography key={"lutronAutomationSelectionTypography" + selectedDevice.id} className={classes.columnHeading}>Automation Devices</Typography>
                    {automatedDevices.map(device =>
                        <Chip key={"automation_chip_" + device} label={device} className={classes.chip} onDelete={() => removeSelectedAutomation(device)} />
                    )}
                </div>
            </ExpansionPanelDetails>
        </ExpansionPanel>
    )
}
