import { useQuery } from '@apollo/client';
import { makeStyles } from '@material-ui/core/styles';
import React from "react";
import { HomeAssistantIntegration } from '../../components/IntegrationArea/HomeAssistant/HomeAssistant';
import { IntegrationArea } from '../../components/IntegrationArea/IntegrationArea';
import { getIntegrations } from './queries/getIntegrations';


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
}));



export function Integration() {
  const { loading, error, data } = useQuery(getIntegrations);

  const classes = useStyles();

  if (loading) return <p>loading ...</p>;
  if (error) {
    console.log(error);
    return <p>Error :( {error}</p>;
  }



  return (
    <div className={classes.root}>
      {data.integrations.integrations.map(integration => {
        var integrationValue = {}
        switch (integration) {
          case "hass":
            integrationValue.name = "Home Assistant";
            integrationValue.value = data.integrations.hass;
            return <HomeAssistantIntegration integration={integrationValue} />
          case "kasa":
            integrationValue.name = "Kasa Smart Home";
            integrationValue.value = data.integrations.kasa;
            return <IntegrationArea integration={integrationValue} />
          case "lutron":
            integrationValue.name = "Lutron Lighting";
            integrationValue.value = data.integrations.lutron;
            return <IntegrationArea integration={integrationValue} />
          default:
            return <div />
        }
      })}
    </div>
  )
}