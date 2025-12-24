import React, { Suspense, useState, useMemo} from "react";
import { Route, Routes } from "react-router-dom";

import { ClockLoader  } from "react-spinners"

import {UserContext, Layout} from "./components";

// Error Page
const Page404 = React.lazy(() => import('./pages/notFound'));
const Table = React.lazy(() => import('./pages/userTable'));
const Detail = React.lazy(() => import('./pages/detail'));
const About = React.lazy(() => import('./pages/about'));

const App = () => {
  const defaultBookContext = {
    users: [],
    page: 1,
    page_count: 10,
    page_size: 20
  };

  const [users, setUsers] = useState(defaultBookContext);
  const contextMemo = useMemo(
    () => ({ users, setUsers }),
    [users]
  );

  return (
    <UserContext.Provider value={contextMemo}>
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
    </UserContext.Provider>
  );
};

export default App;