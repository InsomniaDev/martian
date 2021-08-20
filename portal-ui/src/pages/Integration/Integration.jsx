import { useQuery } from '@apollo/client';
import { makeStyles } from '@material-ui/core/styles';
import React from "react";
import { HomeAssistantIntegration } from '../../components/IntegrationArea/HomeAssistant/HomeAssistant';
import { IntegrationArea } from '../../components/IntegrationArea/IntegrationArea';
import { areaNames } from '../../components/IntegrationArea/queries/areaNames';
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
    padding: `${theme.spacing()}px ${theme.spacing(2)}px`,
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
  const { loading, error, data, refetch } = useQuery(getIntegrations);
  const { loading: areaLoading, error: areaError, data: areaData } = useQuery(areaNames);

  const classes = useStyles();

  if (loading) return <p>loading ...</p>;
  if (error) {
    console.log(error);
    return <p>Error :( {error}</p>;
  }

  if (areaLoading) return <p>loading ...</p>;
  if (areaError) {
    console.log(areaError);
    return <p>Error :( {areaError}</p>;
  }

  return (
    <div className={classes.root}>
      {data.integrations.integrations.map(integration => {
        var integrationValue = {}
        switch (integration) {
          case "hass":
            integrationValue.name = "Home Assistant";
            integrationValue.value = data.integrations.hass;
            return <HomeAssistantIntegration key="hassIntegration" areaData={areaData} integration={integrationValue} refetchData={() => refetch()} />
          case "kasa":
            integrationValue.name = "Kasa Smart Home";
            integrationValue.value = data.integrations.kasa;
            return <IntegrationArea key="kasaIntegration" integration={integrationValue} />
          case "lutron":
            integrationValue.name = "Lutron Lighting";
            integrationValue.value = data.integrations.lutron;
            return <IntegrationArea key="lutronIntegration" integration={integrationValue} />
          default:
            return <div />
        }
      })}
    </div>
  )
}
