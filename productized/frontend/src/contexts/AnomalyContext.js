import React, { createContext, useContext, useReducer } from 'react';
import { anomalyAPI } from '../api';

const AnomalyContext = createContext();

const initialState = {
  anomalies: [],
  currentAnomaly: null,
  loading: false,
  error: null,
};

const anomalyReducer = (state, action) => {
  switch (action.type) {
    case 'FETCH_ANOMALIES_START':
      return {
        ...state,
        loading: true,
        error: null,
      };
    case 'FETCH_ANOMALIES_SUCCESS':
      return {
        ...state,
        anomalies: action.payload,
        loading: false,
      };
    case 'FETCH_ANOMALIES_ERROR':
      return {
        ...state,
        loading: false,
        error: action.payload,
      };
    case 'FETCH_ANOMALY_START':
      return {
        ...state,
        loading: true,
        error: null,
      };
    case 'FETCH_ANOMALY_SUCCESS':
      return {
        ...state,
        currentAnomaly: action.payload,
        loading: false,
      };
    case 'FETCH_ANOMALY_ERROR':
      return {
        ...state,
        loading: false,
        error: action.payload,
      };
    case 'RESOLVE_ANOMALY':
      return {
        ...state,
        anomalies: state.anomalies.map(anomaly =>
          anomaly.id === action.payload.id ? { ...anomaly, status: 'resolved' } : anomaly
        ),
        currentAnomaly: state.currentAnomaly?.id === action.payload.id 
          ? { ...state.currentAnomaly, status: 'resolved' } 
          : state.currentAnomaly,
      };
    case 'CLEAR_CURRENT_ANOMALY':
      return {
        ...state,
        currentAnomaly: null,
      };
    default:
      return state;
  }
};

export const AnomalyProvider = ({ children }) => {
  const [state, dispatch] = useReducer(anomalyReducer, initialState);

  const fetchAnomalies = async (params = {}) => {
    dispatch({ type: 'FETCH_ANOMALIES_START' });
    try {
      const response = await anomalyAPI.getAnomalies(params);
      dispatch({
        type: 'FETCH_ANOMALIES_SUCCESS',
        payload: response.data,
      });
    } catch (error) {
      dispatch({
        type: 'FETCH_ANOMALIES_ERROR',
        payload: error.message,
      });
    }
  };

  const fetchAnomaly = async (id) => {
    dispatch({ type: 'FETCH_ANOMALY_START' });
    try {
      const response = await anomalyAPI.getAnomaly(id);
      dispatch({
        type: 'FETCH_ANOMALY_SUCCESS',
        payload: response.data,
      });
    } catch (error) {
      dispatch({
        type: 'FETCH_ANOMALY_ERROR',
        payload: error.message,
      });
    }
  };

  const resolveAnomaly = async (id, resolution) => {
    try {
      await anomalyAPI.resolveAnomaly(id, resolution);
      dispatch({
        type: 'RESOLVE_ANOMALY',
        payload: { id },
      });
    } catch (error) {
      console.error('Failed to resolve anomaly:', error);
    }
  };

  const clearCurrentAnomaly = () => {
    dispatch({ type: 'CLEAR_CURRENT_ANOMALY' });
  };

  return (
    <AnomalyContext.Provider
      value={{
        ...state,
        fetchAnomalies,
        fetchAnomaly,
        resolveAnomaly,
        clearCurrentAnomaly,
      }}
    >
      {children}
    </AnomalyContext.Provider>
  );
};

export const useAnomaly = () => {
  const context = useContext(AnomalyContext);
  if (!context) {
    throw new Error('useAnomaly must be used within an AnomalyProvider');
  }
  return context;
};