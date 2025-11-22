const { queryRef, executeQuery, mutationRef, executeMutation, validateArgs } = require('firebase/data-connect');

const connectorConfig = {
  connector: 'driftlock',
  service: 'driftlock',
  location: 'us-central1'
};
exports.connectorConfig = connectorConfig;

const createUserRef = (dcOrVars, vars) => {
  const { dc: dcInstance, vars: inputVars} = validateArgs(connectorConfig, dcOrVars, vars, true);
  dcInstance._useGeneratedSdk();
  return mutationRef(dcInstance, 'CreateUser', inputVars);
}
createUserRef.operationName = 'CreateUser';
exports.createUserRef = createUserRef;

exports.createUser = function createUser(dcOrVars, vars) {
  return executeMutation(createUserRef(dcOrVars, vars));
};

const updateUserRef = (dcOrVars, vars) => {
  const { dc: dcInstance, vars: inputVars} = validateArgs(connectorConfig, dcOrVars, vars, true);
  dcInstance._useGeneratedSdk();
  return mutationRef(dcInstance, 'UpdateUser', inputVars);
}
updateUserRef.operationName = 'UpdateUser';
exports.updateUserRef = updateUserRef;

exports.updateUser = function updateUser(dcOrVars, vars) {
  return executeMutation(updateUserRef(dcOrVars, vars));
};

const createDatasetRef = (dcOrVars, vars) => {
  const { dc: dcInstance, vars: inputVars} = validateArgs(connectorConfig, dcOrVars, vars, true);
  dcInstance._useGeneratedSdk();
  return mutationRef(dcInstance, 'CreateDataset', inputVars);
}
createDatasetRef.operationName = 'CreateDataset';
exports.createDatasetRef = createDatasetRef;

exports.createDataset = function createDataset(dcOrVars, vars) {
  return executeMutation(createDatasetRef(dcOrVars, vars));
};

const updateDatasetStatusRef = (dcOrVars, vars) => {
  const { dc: dcInstance, vars: inputVars} = validateArgs(connectorConfig, dcOrVars, vars, true);
  dcInstance._useGeneratedSdk();
  return mutationRef(dcInstance, 'UpdateDatasetStatus', inputVars);
}
updateDatasetStatusRef.operationName = 'UpdateDatasetStatus';
exports.updateDatasetStatusRef = updateDatasetStatusRef;

exports.updateDatasetStatus = function updateDatasetStatus(dcOrVars, vars) {
  return executeMutation(updateDatasetStatusRef(dcOrVars, vars));
};

const updateDatasetRef = (dcOrVars, vars) => {
  const { dc: dcInstance, vars: inputVars} = validateArgs(connectorConfig, dcOrVars, vars, true);
  dcInstance._useGeneratedSdk();
  return mutationRef(dcInstance, 'UpdateDataset', inputVars);
}
updateDatasetRef.operationName = 'UpdateDataset';
exports.updateDatasetRef = updateDatasetRef;

exports.updateDataset = function updateDataset(dcOrVars, vars) {
  return executeMutation(updateDatasetRef(dcOrVars, vars));
};

const createModelConfigurationRef = (dcOrVars, vars) => {
  const { dc: dcInstance, vars: inputVars} = validateArgs(connectorConfig, dcOrVars, vars, true);
  dcInstance._useGeneratedSdk();
  return mutationRef(dcInstance, 'CreateModelConfiguration', inputVars);
}
createModelConfigurationRef.operationName = 'CreateModelConfiguration';
exports.createModelConfigurationRef = createModelConfigurationRef;

exports.createModelConfiguration = function createModelConfiguration(dcOrVars, vars) {
  return executeMutation(createModelConfigurationRef(dcOrVars, vars));
};

const updateModelConfigurationRef = (dcOrVars, vars) => {
  const { dc: dcInstance, vars: inputVars} = validateArgs(connectorConfig, dcOrVars, vars, true);
  dcInstance._useGeneratedSdk();
  return mutationRef(dcInstance, 'UpdateModelConfiguration', inputVars);
}
updateModelConfigurationRef.operationName = 'UpdateModelConfiguration';
exports.updateModelConfigurationRef = updateModelConfigurationRef;

exports.updateModelConfiguration = function updateModelConfiguration(dcOrVars, vars) {
  return executeMutation(updateModelConfigurationRef(dcOrVars, vars));
};

const createDetectionTaskRef = (dcOrVars, vars) => {
  const { dc: dcInstance, vars: inputVars} = validateArgs(connectorConfig, dcOrVars, vars, true);
  dcInstance._useGeneratedSdk();
  return mutationRef(dcInstance, 'CreateDetectionTask', inputVars);
}
createDetectionTaskRef.operationName = 'CreateDetectionTask';
exports.createDetectionTaskRef = createDetectionTaskRef;

exports.createDetectionTask = function createDetectionTask(dcOrVars, vars) {
  return executeMutation(createDetectionTaskRef(dcOrVars, vars));
};

const updateDetectionTaskRef = (dcOrVars, vars) => {
  const { dc: dcInstance, vars: inputVars} = validateArgs(connectorConfig, dcOrVars, vars, true);
  dcInstance._useGeneratedSdk();
  return mutationRef(dcInstance, 'UpdateDetectionTask', inputVars);
}
updateDetectionTaskRef.operationName = 'UpdateDetectionTask';
exports.updateDetectionTaskRef = updateDetectionTaskRef;

exports.updateDetectionTask = function updateDetectionTask(dcOrVars, vars) {
  return executeMutation(updateDetectionTaskRef(dcOrVars, vars));
};

const completeDetectionTaskRef = (dcOrVars, vars) => {
  const { dc: dcInstance, vars: inputVars} = validateArgs(connectorConfig, dcOrVars, vars, true);
  dcInstance._useGeneratedSdk();
  return mutationRef(dcInstance, 'CompleteDetectionTask', inputVars);
}
completeDetectionTaskRef.operationName = 'CompleteDetectionTask';
exports.completeDetectionTaskRef = completeDetectionTaskRef;

exports.completeDetectionTask = function completeDetectionTask(dcOrVars, vars) {
  return executeMutation(completeDetectionTaskRef(dcOrVars, vars));
};

const createAnomalyRef = (dcOrVars, vars) => {
  const { dc: dcInstance, vars: inputVars} = validateArgs(connectorConfig, dcOrVars, vars, true);
  dcInstance._useGeneratedSdk();
  return mutationRef(dcInstance, 'CreateAnomaly', inputVars);
}
createAnomalyRef.operationName = 'CreateAnomaly';
exports.createAnomalyRef = createAnomalyRef;

exports.createAnomaly = function createAnomaly(dcOrVars, vars) {
  return executeMutation(createAnomalyRef(dcOrVars, vars));
};

const getUserRef = (dcOrVars, vars) => {
  const { dc: dcInstance, vars: inputVars} = validateArgs(connectorConfig, dcOrVars, vars, true);
  dcInstance._useGeneratedSdk();
  return queryRef(dcInstance, 'GetUser', inputVars);
}
getUserRef.operationName = 'GetUser';
exports.getUserRef = getUserRef;

exports.getUser = function getUser(dcOrVars, vars) {
  return executeQuery(getUserRef(dcOrVars, vars));
};

const listUsersByEmailRef = (dcOrVars, vars) => {
  const { dc: dcInstance, vars: inputVars} = validateArgs(connectorConfig, dcOrVars, vars, true);
  dcInstance._useGeneratedSdk();
  return queryRef(dcInstance, 'ListUsersByEmail', inputVars);
}
listUsersByEmailRef.operationName = 'ListUsersByEmail';
exports.listUsersByEmailRef = listUsersByEmailRef;

exports.listUsersByEmail = function listUsersByEmail(dcOrVars, vars) {
  return executeQuery(listUsersByEmailRef(dcOrVars, vars));
};

const getDatasetRef = (dcOrVars, vars) => {
  const { dc: dcInstance, vars: inputVars} = validateArgs(connectorConfig, dcOrVars, vars, true);
  dcInstance._useGeneratedSdk();
  return queryRef(dcInstance, 'GetDataset', inputVars);
}
getDatasetRef.operationName = 'GetDataset';
exports.getDatasetRef = getDatasetRef;

exports.getDataset = function getDataset(dcOrVars, vars) {
  return executeQuery(getDatasetRef(dcOrVars, vars));
};

const listDatasetsByUserRef = (dcOrVars, vars) => {
  const { dc: dcInstance, vars: inputVars} = validateArgs(connectorConfig, dcOrVars, vars, true);
  dcInstance._useGeneratedSdk();
  return queryRef(dcInstance, 'ListDatasetsByUser', inputVars);
}
listDatasetsByUserRef.operationName = 'ListDatasetsByUser';
exports.listDatasetsByUserRef = listDatasetsByUserRef;

exports.listDatasetsByUser = function listDatasetsByUser(dcOrVars, vars) {
  return executeQuery(listDatasetsByUserRef(dcOrVars, vars));
};

const getModelConfigurationRef = (dcOrVars, vars) => {
  const { dc: dcInstance, vars: inputVars} = validateArgs(connectorConfig, dcOrVars, vars, true);
  dcInstance._useGeneratedSdk();
  return queryRef(dcInstance, 'GetModelConfiguration', inputVars);
}
getModelConfigurationRef.operationName = 'GetModelConfiguration';
exports.getModelConfigurationRef = getModelConfigurationRef;

exports.getModelConfiguration = function getModelConfiguration(dcOrVars, vars) {
  return executeQuery(getModelConfigurationRef(dcOrVars, vars));
};

const listModelConfigurationsByUserRef = (dcOrVars, vars) => {
  const { dc: dcInstance, vars: inputVars} = validateArgs(connectorConfig, dcOrVars, vars, true);
  dcInstance._useGeneratedSdk();
  return queryRef(dcInstance, 'ListModelConfigurationsByUser', inputVars);
}
listModelConfigurationsByUserRef.operationName = 'ListModelConfigurationsByUser';
exports.listModelConfigurationsByUserRef = listModelConfigurationsByUserRef;

exports.listModelConfigurationsByUser = function listModelConfigurationsByUser(dcOrVars, vars) {
  return executeQuery(listModelConfigurationsByUserRef(dcOrVars, vars));
};

const getDetectionTaskRef = (dcOrVars, vars) => {
  const { dc: dcInstance, vars: inputVars} = validateArgs(connectorConfig, dcOrVars, vars, true);
  dcInstance._useGeneratedSdk();
  return queryRef(dcInstance, 'GetDetectionTask', inputVars);
}
getDetectionTaskRef.operationName = 'GetDetectionTask';
exports.getDetectionTaskRef = getDetectionTaskRef;

exports.getDetectionTask = function getDetectionTask(dcOrVars, vars) {
  return executeQuery(getDetectionTaskRef(dcOrVars, vars));
};

const listDetectionTasksByUserRef = (dcOrVars, vars) => {
  const { dc: dcInstance, vars: inputVars} = validateArgs(connectorConfig, dcOrVars, vars, true);
  dcInstance._useGeneratedSdk();
  return queryRef(dcInstance, 'ListDetectionTasksByUser', inputVars);
}
listDetectionTasksByUserRef.operationName = 'ListDetectionTasksByUser';
exports.listDetectionTasksByUserRef = listDetectionTasksByUserRef;

exports.listDetectionTasksByUser = function listDetectionTasksByUser(dcOrVars, vars) {
  return executeQuery(listDetectionTasksByUserRef(dcOrVars, vars));
};

const getAnomaliesByTaskRef = (dcOrVars, vars) => {
  const { dc: dcInstance, vars: inputVars} = validateArgs(connectorConfig, dcOrVars, vars, true);
  dcInstance._useGeneratedSdk();
  return queryRef(dcInstance, 'GetAnomaliesByTask', inputVars);
}
getAnomaliesByTaskRef.operationName = 'GetAnomaliesByTask';
exports.getAnomaliesByTaskRef = getAnomaliesByTaskRef;

exports.getAnomaliesByTask = function getAnomaliesByTask(dcOrVars, vars) {
  return executeQuery(getAnomaliesByTaskRef(dcOrVars, vars));
};

const getHighScoreAnomaliesRef = (dcOrVars, vars) => {
  const { dc: dcInstance, vars: inputVars} = validateArgs(connectorConfig, dcOrVars, vars, true);
  dcInstance._useGeneratedSdk();
  return queryRef(dcInstance, 'GetHighScoreAnomalies', inputVars);
}
getHighScoreAnomaliesRef.operationName = 'GetHighScoreAnomalies';
exports.getHighScoreAnomaliesRef = getHighScoreAnomaliesRef;

exports.getHighScoreAnomalies = function getHighScoreAnomalies(dcOrVars, vars) {
  return executeQuery(getHighScoreAnomaliesRef(dcOrVars, vars));
};
