import React, { Suspense } from "react";
import { Route, Routes } from "react-router-dom";

import { Layout, Loading } from "@backpack";

// Error Page
const Page404 = React.lazy(() => import('@pages/notFound'));
const Table = React.lazy(() => import('@pages/userTable'));
const Detail = React.lazy(() => import('@pages/detail'));
const About = React.lazy(() => import('@pages/about'));

const App = () =>  (
  <Layout>
    <Suspense fallback={<Loading />}>
      <Routes>
        <Route index path="/" element={<Table />} />
        <Route path="/about" element={<About />} />
        <Route path="/detail/:id?" element={<Detail />} />
        <Route path="*" element={<Page404 />} />
      </Routes>
    </Suspense>
  </Layout>
);

export default App;