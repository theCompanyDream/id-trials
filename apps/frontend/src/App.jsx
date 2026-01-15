import React, { Suspense } from "react";
import { Route, Routes } from "react-router-dom";

import { ClockLoader  } from "react-spinners"

import { Layout} from "./components";

// Error Page
const Page404 = React.lazy(() => import('./pages/notFound'));
const Table = React.lazy(() => import('./pages/userTable'));
const Detail = React.lazy(() => import('./pages/detail'));
const About = React.lazy(() => import('./pages/about'));

const App = () =>  (
  <Layout>
    <Suspense fallback={<ClockLoader  color="#FFF200" size={50} />}>
      <Routes>
        <Route index path="/" element={<Table />} />
        <Route path="/detail/:id?" element={<Detail />} />
        <Route path="*" element={<Page404 />} />
        <Route path="/about" element={<About />} />
      </Routes>
    </Suspense>
  </Layout>
);

export default App;