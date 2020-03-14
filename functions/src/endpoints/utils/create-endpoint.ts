type FirebaseApp = import('firebase-admin').app.App;

type EndpointFunction = import('firebase-functions').HttpsFunction;

export function createEndpoint(factory: (app: FirebaseApp) => EndpointFunction) {
  return (app: FirebaseApp) => factory(app);
}
