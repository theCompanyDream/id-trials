import React, { Suspense } from "react";
import { Route, Routes } from "react-router-dom";

import { Layout, Loading } from "@backpack";

// Error Page
const About = React.lazy(() => import('@pages/about'));
const Analytics = React.lazy(() => import('@pages/analysis'))
const Detail = React.lazy(() => import('@pages/detail'));
const Page404 = React.lazy(() => import('@pages/notFound'));
const Table = React.lazy(() => import('@pages/userTable'));

const App = () =>  (
  <Layout>
    <Suspense fallback={<Loading />}>
      <Routes>
        <Route index path="/" element={<Analytics />} />
        <Route path="/about" element={<About />} />
        <Route path="/detail/:id?" element={<Detail />} />
        <Route path="/explore" element={<Table />} />
        <Route path="*" element={<Page404 />} />
      </Routes>
    </Suspense>
  </Layout>
);

export default App;