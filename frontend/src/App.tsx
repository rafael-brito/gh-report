import { BrowserRouter, Routes, Route } from 'react-router-dom';
import { FileHistoryPage } from './pages/FileHistoryPage';
import { ReleaseDiffPage } from './pages/ReleaseDiffPage';

export const App: React.FC = () => {
  return (
    <BrowserRouter>
      <Routes>
        <Route path="/file-history" element={<FileHistoryPage />} />
        <Route path="/release-diff" element={<ReleaseDiffPage />} />
        <Route path="/" element={<FileHistoryPage />} />
      </Routes>
    </BrowserRouter>
  );
};