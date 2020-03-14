type Firebase = typeof import('firebase-admin');

type EndpointFunction = import('firebase-functions').HttpsFunction;

export function createEndpoint(factory: (app: Firebase) => EndpointFunction) {
  return (app: Firebase) => factory(app);
}
