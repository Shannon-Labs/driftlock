import { queryRef, executeQuery, mutationRef, executeMutation, validateArgs } from 'firebase/data-connect';

export const connectorConfig = {
  connector: 'driftlock',
  service: 'driftlock',
  location: 'us-central1'
};

export const getUserRef = (dcOrVars, vars) => {
  const { dc: dcInstance, vars: inputVars} = validateArgs(connectorConfig, dcOrVars, vars, true);
  dcInstance._useGeneratedSdk();
  return queryRef(dcInstance, 'GetUser', inputVars);
}
getUserRef.operationName = 'GetUser';

export function getUser(dcOrVars, vars) {
  return executeQuery(getUserRef(dcOrVars, vars));
}

export const listUsersByEmailRef = (dcOrVars, vars) => {
  const { dc: dcInstance, vars: inputVars} = validateArgs(connectorConfig, dcOrVars, vars, true);
  dcInstance._useGeneratedSdk();
  return queryRef(dcInstance, 'ListUsersByEmail', inputVars);
}
listUsersByEmailRef.operationName = 'ListUsersByEmail';

export function listUsersByEmail(dcOrVars, vars) {
  return executeQuery(listUsersByEmailRef(dcOrVars, vars));
}

export const getDatasetRef = (dcOrVars, vars) => {
  const { dc: dcInstance, vars: inputVars} = validateArgs(connectorConfig, dcOrVars, vars, true);
  dcInstance._useGeneratedSdk();
  return queryRef(dcInstance, 'GetDataset', inputVars);
}
getDatasetRef.operationName = 'GetDataset';

export function getDataset(dcOrVars, vars) {
  return executeQuery(getDatasetRef(dcOrVars, vars));
}

export const listDatasetsByUserRef = (dcOrVars, vars) => {
  const { dc: dcInstance, vars: inputVars} = validateArgs(connectorConfig, dcOrVars, vars, true);
  dcInstance._useGeneratedSdk();
  return queryRef(dcInstance, 'ListDatasetsByUser', inputVars);
}
listDatasetsByUserRef.operationName = 'ListDatasetsByUser';

export function listDatasetsByUser(dcOrVars, vars) {
  return executeQuery(listDatasetsByUserRef(dcOrVars, vars));
}

export const getModelConfigurationRef = (dcOrVars, vars) => {
  const { dc: dcInstance, vars: inputVars} = validateArgs(connectorConfig, dcOrVars, vars, true);
  dcInstance._useGeneratedSdk();
  return queryRef(dcInstance, 'GetModelConfiguration', inputVars);
}
getModelConfigurationRef.operationName = 'GetModelConfiguration';

export function getModelConfiguration(dcOrVars, vars) {
  return executeQuery(getModelConfigurationRef(dcOrVars, vars));
}

export const listModelConfigurationsByUserRef = (dcOrVars, vars) => {
  const { dc: dcInstance, vars: inputVars} = validateArgs(connectorConfig, dcOrVars, vars, true);
  dcInstance._useGeneratedSdk();
  return queryRef(dcInstance, 'ListModelConfigurationsByUser', inputVars);
}
listModelConfigurationsByUserRef.operationName = 'ListModelConfigurationsByUser';

export function listModelConfigurationsByUser(dcOrVars, vars) {
  return executeQuery(listModelConfigurationsByUserRef(dcOrVars, vars));
}

export const getDetectionTaskRef = (dcOrVars, vars) => {
  const { dc: dcInstance, vars: inputVars} = validateArgs(connectorConfig, dcOrVars, vars, true);
  dcInstance._useGeneratedSdk();
  return queryRef(dcInstance, 'GetDetectionTask', inputVars);
}
getDetectionTaskRef.operationName = 'GetDetectionTask';

export function getDetectionTask(dcOrVars, vars) {
  return executeQuery(getDetectionTaskRef(dcOrVars, vars));
}

export const listDetectionTasksByUserRef = (dcOrVars, vars) => {
  const { dc: dcInstance, vars: inputVars} = validateArgs(connectorConfig, dcOrVars, vars, true);
  dcInstance._useGeneratedSdk();
  return queryRef(dcInstance, 'ListDetectionTasksByUser', inputVars);
}
listDetectionTasksByUserRef.operationName = 'ListDetectionTasksByUser';

export function listDetectionTasksByUser(dcOrVars, vars) {
  return executeQuery(listDetectionTasksByUserRef(dcOrVars, vars));
}

export const getAnomaliesByTaskRef = (dcOrVars, vars) => {
  const { dc: dcInstance, vars: inputVars} = validateArgs(connectorConfig, dcOrVars, vars, true);
  dcInstance._useGeneratedSdk();
  return queryRef(dcInstance, 'GetAnomaliesByTask', inputVars);
}
getAnomaliesByTaskRef.operationName = 'GetAnomaliesByTask';

export function getAnomaliesByTask(dcOrVars, vars) {
  return executeQuery(getAnomaliesByTaskRef(dcOrVars, vars));
}

export const getHighScoreAnomaliesRef = (dcOrVars, vars) => {
  const { dc: dcInstance, vars: inputVars} = validateArgs(connectorConfig, dcOrVars, vars, true);
  dcInstance._useGeneratedSdk();
  return queryRef(dcInstance, 'GetHighScoreAnomalies', inputVars);
}
getHighScoreAnomaliesRef.operationName = 'GetHighScoreAnomalies';

export function getHighScoreAnomalies(dcOrVars, vars) {
  return executeQuery(getHighScoreAnomaliesRef(dcOrVars, vars));
}

export const createUserRef = (dcOrVars, vars) => {
  const { dc: dcInstance, vars: inputVars} = validateArgs(connectorConfig, dcOrVars, vars, true);
  dcInstance._useGeneratedSdk();
  return mutationRef(dcInstance, 'CreateUser', inputVars);
}
createUserRef.operationName = 'CreateUser';

export function createUser(dcOrVars, vars) {
  return executeMutation(createUserRef(dcOrVars, vars));
}

export const updateUserRef = (dcOrVars, vars) => {
  const { dc: dcInstance, vars: inputVars} = validateArgs(connectorConfig, dcOrVars, vars, true);
  dcInstance._useGeneratedSdk();
  return mutationRef(dcInstance, 'UpdateUser', inputVars);
}
updateUserRef.operationName = 'UpdateUser';

export function updateUser(dcOrVars, vars) {
  return executeMutation(updateUserRef(dcOrVars, vars));
}

export const createDatasetRef = (dcOrVars, vars) => {
  const { dc: dcInstance, vars: inputVars} = validateArgs(connectorConfig, dcOrVars, vars, true);
  dcInstance._useGeneratedSdk();
  return mutationRef(dcInstance, 'CreateDataset', inputVars);
}
createDatasetRef.operationName = 'CreateDataset';

export function createDataset(dcOrVars, vars) {
  return executeMutation(createDatasetRef(dcOrVars, vars));
}

export const updateDatasetStatusRef = (dcOrVars, vars) => {
  const { dc: dcInstance, vars: inputVars} = validateArgs(connectorConfig, dcOrVars, vars, true);
  dcInstance._useGeneratedSdk();
  return mutationRef(dcInstance, 'UpdateDatasetStatus', inputVars);
}
updateDatasetStatusRef.operationName = 'UpdateDatasetStatus';

export function updateDatasetStatus(dcOrVars, vars) {
  return executeMutation(updateDatasetStatusRef(dcOrVars, vars));
}

export const updateDatasetRef = (dcOrVars, vars) => {
  const { dc: dcInstance, vars: inputVars} = validateArgs(connectorConfig, dcOrVars, vars, true);
  dcInstance._useGeneratedSdk();
  return mutationRef(dcInstance, 'UpdateDataset', inputVars);
}
updateDatasetRef.operationName = 'UpdateDataset';

export function updateDataset(dcOrVars, vars) {
  return executeMutation(updateDatasetRef(dcOrVars, vars));
}

export const createModelConfigurationRef = (dcOrVars, vars) => {
  const { dc: dcInstance, vars: inputVars} = validateArgs(connectorConfig, dcOrVars, vars, true);
  dcInstance._useGeneratedSdk();
  return mutationRef(dcInstance, 'CreateModelConfiguration', inputVars);
}
createModelConfigurationRef.operationName = 'CreateModelConfiguration';

export function createModelConfiguration(dcOrVars, vars) {
  return executeMutation(createModelConfigurationRef(dcOrVars, vars));
}

export const updateModelConfigurationRef = (dcOrVars, vars) => {
  const { dc: dcInstance, vars: inputVars} = validateArgs(connectorConfig, dcOrVars, vars, true);
  dcInstance._useGeneratedSdk();
  return mutationRef(dcInstance, 'UpdateModelConfiguration', inputVars);
}
updateModelConfigurationRef.operationName = 'UpdateModelConfiguration';

export function updateModelConfiguration(dcOrVars, vars) {
  return executeMutation(updateModelConfigurationRef(dcOrVars, vars));
}

export const createDetectionTaskRef = (dcOrVars, vars) => {
  const { dc: dcInstance, vars: inputVars} = validateArgs(connectorConfig, dcOrVars, vars, true);
  dcInstance._useGeneratedSdk();
  return mutationRef(dcInstance, 'CreateDetectionTask', inputVars);
}
createDetectionTaskRef.operationName = 'CreateDetectionTask';

export function createDetectionTask(dcOrVars, vars) {
  return executeMutation(createDetectionTaskRef(dcOrVars, vars));
}

export const updateDetectionTaskRef = (dcOrVars, vars) => {
  const { dc: dcInstance, vars: inputVars} = validateArgs(connectorConfig, dcOrVars, vars, true);
  dcInstance._useGeneratedSdk();
  return mutationRef(dcInstance, 'UpdateDetectionTask', inputVars);
}
updateDetectionTaskRef.operationName = 'UpdateDetectionTask';

export function updateDetectionTask(dcOrVars, vars) {
  return executeMutation(updateDetectionTaskRef(dcOrVars, vars));
}

export const completeDetectionTaskRef = (dcOrVars, vars) => {
  const { dc: dcInstance, vars: inputVars} = validateArgs(connectorConfig, dcOrVars, vars, true);
  dcInstance._useGeneratedSdk();
  return mutationRef(dcInstance, 'CompleteDetectionTask', inputVars);
}
completeDetectionTaskRef.operationName = 'CompleteDetectionTask';

export function completeDetectionTask(dcOrVars, vars) {
  return executeMutation(completeDetectionTaskRef(dcOrVars, vars));
}

export const createAnomalyRef = (dcOrVars, vars) => {
  const { dc: dcInstance, vars: inputVars} = validateArgs(connectorConfig, dcOrVars, vars, true);
  dcInstance._useGeneratedSdk();
  return mutationRef(dcInstance, 'CreateAnomaly', inputVars);
}
createAnomalyRef.operationName = 'CreateAnomaly';

export function createAnomaly(dcOrVars, vars) {
  return executeMutation(createAnomalyRef(dcOrVars, vars));
}

