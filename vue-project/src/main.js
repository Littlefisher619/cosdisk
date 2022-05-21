import App from './App.vue'

import { ApolloClient, InMemoryCache, createHttpLink, ApolloLink } from '@apollo/client/core'
import { onError } from "apollo-link-error"
import Antd from 'ant-design-vue';
import 'ant-design-vue/dist/antd.css';

const cache = new InMemoryCache()

const authLink = new ApolloLink((operation, forward) => {
  operation.setContext(({ headers }) => ({
    headers: {
      ...headers,
      authorization: localStorage.getItem('token') ? `Bearer ${localStorage.getItem('token')}` : null,
    },
  }));
  return forward(operation); // Go to the next link in the chain. Similar to `next` in Express.js middleware.
})


// setContext(async (_, { headers }) => {
//   // Use your async token function here:
//   // Return the headers to the context so httpLink can read them
//   if (localStorage.getItem('Authentication')) {
//     return {
//       headers: {
//         ...headers,
//         authorization: 'Bearer ' + localStorage.getItem('Authentication')
//       }
//     }
//   } else {
//     return {
//       headers: {
//         ...headers,
//       }
//     }
//   }
// })

import { createUploadLink } from "apollo-upload-client";


const API_URL = process.env.VUE_APP_GQL_BACKEND ?? "/query"
const uploadlink = createUploadLink({ uri: API_URL });

import { message as msg } from 'ant-design-vue';
import 'ant-design-vue/es/message/style/css';

const errorLink = onError(({ graphQLErrors, networkError }) => {
  if (graphQLErrors) {
    msg.destroy()
    graphQLErrors.map(({ message, locations, path }) => {
      if (message.includes("Access denied")) {
        msg.error('登录态失效或尚未登录，请重新登录！')
        localStorage.clear()
        setTimeout(() => {
          router.replace('/login')
        }, 1000)
      } else {
        msg.error('API错误：'+ message)
      }
      console.log(
        `[GraphQL error]: Message: ${message}, Location: ${locations}, Path: ${path}`
      )
    })


  }
  if (networkError) {
    console.log(`[Network error]: ${networkError}`)
    msg.error('网络错误: ' + networkError.message)
  }

})

const links = [errorLink, authLink, uploadlink]
// Create the apollo client
export const apolloClient = new ApolloClient({
  link: ApolloLink.from(links),
  uri: API_URL,
  //connectToDevTools: true,
  cache,
})

import { createApolloProvider } from '@vue/apollo-option'
import { ApolloClients } from "@vue/apollo-composable";
import { setContext } from 'apollo-link-context'

const apolloProvider = createApolloProvider({
  defaultClient: apolloClient,
})

import { createApp, h, provide } from 'vue'
import './index.css'
import router from './router/router.js'

const app = createApp({
  setup() {
    provide(ApolloClients, {
      default: apolloClient,
    });
  },
  render: () => h(App),
});

app.use(router)
// app.use(apolloProvider)
app.use(Antd)
app.mount('#app')

