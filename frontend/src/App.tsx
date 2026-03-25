import React from 'react';
import { BrowserRouter, Routes, Route } from 'react-router-dom';
import { AppLayout } from './components/AppLayout';
import { HomePage } from './pages/HomePage';
import { FileHistoryPage } from './pages/FileHistoryPage';
import { ReleaseDiffPage } from './pages/ReleaseDiffPage';

export const App: React.FC = () => {
  return (
    <BrowserRouter>
      <AppLayout>
        <Routes>
          <Route path="/" element={<HomePage />} />
          <Route path="/file-history" element={<FileHistoryPage />} />
          <Route path="/release-diff" element={<ReleaseDiffPage />} />
        </Routes>
      </AppLayout>
    </BrowserRouter>
  );
};