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
import HarmonyEditMenu from './HarmonyEditMenu';


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



export function HarmonyIntegration({ integration, refetchData, areaData }) {
    const classes = useStyles();

    const [selectedDevice, changeSelectedDevice] = useState("");
    const [selectDevicesForIntegrationMutation] = useMutation(selectDevicesForIntegration);
    const handleChange = (event) => {
        changeSelectedDevice(event.target.value);
    };

    const clearSelected = () => {
        changeSelectedDevice("");
    }

    const updateSelected = (newDevice) => {
        changeSelectedDevice(newDevice);
    }

    return (
        <ExpansionPanel key="harmonyExpansionPanel">
            <ExpansionPanelSummary key="harmony" >
                <div key="harmonyTypographyNameDiv" className={classes.column}>
                    <Typography key="harmonyTypographyName" className={classes.heading}>{integration.name}</Typography>
                </div>
                <div key="harmonyTypographyHeadingDiv" className={classes.column}>
                    <Typography key="harmonyTypographyHeading" className={classes.secondaryHeading}>Edit Configuration</Typography>
                </div>
            </ExpansionPanelSummary>
            <ExpansionPanelDetails key="harmonyExpansionPanelDetails" className={classes.details}>
                <div key="harmonyExpansionPanelDetailsDiv" className={classes.column}>
                    <div className={classes.deviceDetails}>
                        <Typography
                            key={"selectedDeviceIpAddress" + integration.value.ipAddress}
                            className={classes.deviceHeading}
                            align="left">
                            <em className={classes.em}>ACTIVITY ID:</em>     {integration.value.activityId}
                        </Typography>
                        <Typography
                            key={"selectedDeviceType" + integration.value.ipAddress}
                            className={classes.deviceHeading}
                            align="left">
                            <em className={classes.em}>IP ADDRESS:</em>          {integration.value.ipAddress.split(":")[0]}
                        </Typography>
                        <Typography
                            key={"selectedDeviceAreaName" + integration.value.ipAddress}
                            className={classes.deviceHeading}
                            align="left">
                            <em className={classes.em}>NAME:</em>     {integration.value.name}
                        </Typography>
                        <Typography
                            key={"selectedDeviceAreaName" + integration.value.ipAddress}
                            className={classes.deviceHeading}
                            align="left">
                            <em className={classes.em}>AREA NAME:</em>     {integration.value.areaName}
                        </Typography>
                        <div key={"selectedDeviceDiv" + integration.value.ipAddress} className={classes.buttonDiv}>
                            <HarmonyEditMenu
                                key={"selectedDeviceEditButton" + integration.value.ipAddress}
                                device={integration.value}
                                buttonStyle={classes.button}
                                buttonText="edit device"
                                areaData={areaData}
                                refetchData={refetchData}
                                updateSelected={updateSelected} />
                        </div>
                    </div>
                </div>
            </ExpansionPanelDetails>
        </ExpansionPanel>
    )
}
