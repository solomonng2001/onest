import { useState, type FormEvent } from "react";
import "./UENPage.css";

type UENResponse = {
  uen: string;
  valid: boolean;
  format?: string;
  message: string;
};

const UEN_API_URL = "http://localhost:8080/api/uen/validate";

function UENPage() {
  const [uen, setUen] = useState("");
  const [result, setResult] = useState<UENResponse | null>(null);
  const [isLoading, setIsLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);

  async function validateUEN(event: FormEvent<HTMLFormElement>) {
    event.preventDefault();

    const trimmedUEN = uen.trim();

    setResult(null);
    setError(null);
    setIsLoading(true);

    try {
      const response = await fetch(UEN_API_URL, {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
        },
        body: JSON.stringify({
          uen: trimmedUEN,
        }),
      });

      if (!response.ok) {
        throw new Error(
          `UEN request failed with status ${response.status}`,
        );
      }

      const responseResult: UENResponse = await response.json();
      setResult(responseResult);
    } catch (requestError) {
      setError(
        requestError instanceof Error
          ? requestError.message
          : "Unable to validate the UEN.",
      );
    } finally {
      setIsLoading(false);
    }
  }

  function handleUENChange(value: string) {
    setUen(value.toUpperCase());
    setResult(null);
    setError(null);
  }

  return (
    <section
      id="uen-panel"
      className="page"
      role="tabpanel"
      aria-labelledby="uen-tab"
    >
      <div className="page-heading">
        <h2>UEN validation</h2>
        <p>Check whether a Unique Entity Number has a valid format.</p>
      </div>

      <div className="uen-card">
        <form className="uen-form" onSubmit={validateUEN}>
          <label htmlFor="uen">Unique Entity Number</label>

          <div className="uen-input-row">
            <input
              id="uen"
              type="text"
              value={uen}
              onChange={(event) => handleUENChange(event.target.value)}
              placeholder="For example, 200912345N"
              autoComplete="off"
              required
            />

            <button
              type="submit"
              disabled={isLoading || uen.trim().length === 0}
            >
              {isLoading ? "Validating..." : "Validate UEN"}
            </button>
          </div>

          <p className="input-hint">
            Enter the UEN without spaces or symbols.
          </p>
        </form>

        {result && (
          <div
            className={
              result.valid
                ? "uen-result uen-result-valid"
                : "uen-result uen-result-invalid"
            }
            role="status"
          >
            <p className="result-title">
              {result.valid ? "Valid UEN format" : "Invalid UEN format"}
            </p>

            <p>{result.message}</p>

            {result.format && (
              <p className="result-format">
                Format:{" "}
                <strong>{result.format.replaceAll("_", " ")}</strong>
              </p>
            )}
          </div>
        )}

        {error && (
          <div className="uen-result uen-result-invalid" role="alert">
            <p className="result-title">Unable to validate UEN</p>
            <p>{error}</p>
          </div>
        )}
      </div>
    </section>
  );
}

export default UENPage;
