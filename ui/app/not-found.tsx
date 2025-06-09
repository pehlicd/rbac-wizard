// ui/app/not-found.tsx
export default function NotFound() {
  return (
    <div className="min-h-screen flex flex-col justify-center items-center bg-gray-100 text-center">
      <h1 className="text-4xl font-bold text-red-600 mb-4">404</h1>
      <p className="text-xl font-medium">Oops! This page could not be found.</p>
      <a href="/" className="mt-6 text-blue-500 hover:underline">Go back to Home</a>
    </div>
  );
}
